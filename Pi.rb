def pi_nilakantha(n)
  sum = 3.0
  i = 2
  sign = 1
  while i <= n
    sum += (sign * 4.0) / (i * (i + 1) * (i + 2))
    i += 2
    sign = -sign
  end
  sum
end
puts pi_nilakantha(2000)
# => 3.141592653340544
