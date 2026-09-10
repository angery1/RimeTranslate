local temp_dir = os.getenv("TEMP") or "."
local bridge_dir = temp_dir .. "\\rime_ollama_bridge"

local request_file  = bridge_dir .. "\\request.txt"
local response_file = bridge_dir .. "\\response.txt"

local function file_exists(path)
    local f = io.open(path, "rb")
    if f then
        f:close()
        return true
    end
    return false
end

local function request_translation(text)
    -- RimeTranslate.exe creates bridge_dir at startup.
    -- IMPORTANT: never call os.execute() here; a Lua filter runs on Rime's
    -- candidate-generation path and synchronous process creation can freeze typing.
    local f = io.open(request_file, "wb")
    if f then
        f:write(text)
        f:close()
        return true
    end
    return false
end

local function read_translation()
    local f = io.open(response_file, "rb")
    if not f then
        return nil, nil
    end

    local source = f:read("*l")
    local translation = f:read("*a")
    f:close()

    if translation then
        translation = translation:gsub("^%s+", ""):gsub("%s+$", "")
    end

    return source, translation
end

local function has_han(s)
    return s:find("[\228-\233]") ~= nil
end

local function has_latin(s)
    return s:find("[A-Za-z]") ~= nil
end

local function is_translatable(s)
    return has_han(s) or has_latin(s)
end

local function init(env)
    env.last_requested = ""

    -- Reset request suppression whenever the Rime translation switch changes,
    -- even if the user toggles it while no composition is active.
    env.option_connection = env.engine.context.option_update_notifier:connect(
        function(_, name)
            if name == "ollama_translation" then
                env.last_requested = ""
            end
        end
    )
end

local function fini(env)
    if env.option_connection then
        env.option_connection:disconnect()
        env.option_connection = nil
    end
end

local function filter(input, env)
    local context = env.engine.context

    -- Translation OFF: pure pass-through. No Ollama request.
    if not context:get_option("ollama_translation") then
        env.last_requested = ""
        for cand in input:iter() do
            yield(cand)
        end
        return
    end

    local first_seen = false

    -- Stream candidates instead of collecting the whole candidate list into a table.
    -- This keeps Rime's candidate path fast even when dictionaries return many items.
    for cand in input:iter() do
        if not first_seen then
            first_seen = true
            local text = cand.text

            if is_translatable(text) then
                local source, translation = read_translation()
                local response_ok = source == text
                    and translation
                    and translation ~= ""
                    and translation ~= text

                -- If the old response was released/removed, allow the same text to
                -- be requested again. request.txt acts as the pending-request marker.
                if not response_ok
                    and (text ~= env.last_requested or not file_exists(request_file)) then
                    if request_translation(text) then
                        env.last_requested = text
                    end
                end

                -- Candidate #1 stays the original text.
                yield(cand)

                -- Insert translation as candidate #2 only when the response belongs
                -- to the current first candidate.
                if response_ok then
                    local translated = Candidate(
                        "ollama_translation",
                        cand.start,
                        cand._end,
                        translation,
                        " 🌐"
                    )
                    yield(translated)
                end
            else
                yield(cand)
            end
        else
            yield(cand)
        end
    end
end

return { init = init, func = filter, fini = fini }
