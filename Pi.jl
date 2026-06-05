function chudnovsky_pi(digits::Int)
    # Set precision slightly higher to avoid rounding errors at the edge
    setprecision(digits * 4 + 100) do
        C = 640320
        C3_OVER_24 = C^3 ÷ 24
        
        # Binary Splitting Function
        function bs(a::Int, b::Int)
            if b - a == 1
                if a == 0
                    return (big(1), big(1), big(0))
                end
                p = (6a - 5) * (2a - 1) * (6a - 1)
                q = a^3 * C3_OVER_24
                t = p * (13591409 + 545140134a)
                if isodd(a)
                    t = -t
                end
                return (big(p), big(q), big(t))
            else
                m = (a + b) ÷ 2
                Pam, Qam, Tam = bs(a, m)
                Pmb, Qmb, Tmb = bs(m, b)
                return (Pam * Pmb, Qam * Qmb, Qmb * Tam + Pam * Tmb)
            end
        end

        # Estimate number of terms needed (14 digits per term)
        n_terms = max(1, ceil(Int, digits / 14.18))
        P, Q, T = bs(0, n_terms)
        
        # Final calculation: pi = (426880 * sqrt(10005) * Q) / (13591409*Q + T) 
        # Note: T in this implementation includes the 'a' term summation logic directly
        # Standard formula adjustment for the specific bs return values:
        numerator = Q * 426880 * sqrt(big(10005))
        denominator = 13591409 * Q + T
        
        return numerator / denominator
    end
end

# Usage:
println(chudnovsky_pi(1000)) # Computes 1000 digits   
