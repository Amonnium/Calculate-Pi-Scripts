function calculatePiLeibniz(iterations)
    local pi = 0
    for i = 0, iterations - 1 do
        local term = 1 / (2 * i + 1)
        if i % 2 == 0 then
            pi = pi + term
        else
            pi = pi - term
        end
    end
    return pi * 4
end

print(calculatePiLeibniz(1000000))      
