local function processor(key, env)
    if key:repr() == "F24" then
        local context = env.engine.context
        if context:is_composing()
            and context:get_option("ollama_translation") then
            context:refresh_non_confirmed_composition()
        end
        return 1
    end
    return 2
end

return { func = processor }
