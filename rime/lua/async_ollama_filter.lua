local temp_dir = os.getenv("TEMP") or "."
local bridge_dir = temp_dir .. "\\rime_ollama_bridge"

local request_file  = bridge_dir .. "\\request.txt"
local response_file = bridge_dir .. "\\response.txt"

local last_requested = ""

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

local function filter(input, env)
    local context = env.engine.context

    -- Translation OFF: pure pass-through. No file I/O, no Ollama request.
    if not context:get_option("ollama_translation") then
        last_requested = ""
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
                if text ~= last_requested then
                    if request_translation(text) then
                        last_requested = text
                    end
                end

                local source, translation = read_translation()

                -- Candidate #1 stays the original text.
                yield(cand)

                -- Insert translation as candidate #2 only when the response belongs
                -- to the current first candidate.
                if source == text
                    and translation
                    and translation ~= ""
                    and translation ~= text then

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

return { func = filter }
