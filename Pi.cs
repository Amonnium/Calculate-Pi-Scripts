using System;
using System.Numerics;
using System.Collections.Generic;
using System.Linq;

public static class BigMath
{
    // Calculate Pi with the specified number of digits of precision using Machin's formula
    public static BigInteger GetPi(int digits, int iterations)
    {
        return 16 * ArcTan1OverX(5, digits).ElementAt(iterations)
             - 4 * ArcTan1OverX(239, digits).ElementAt(iterations);
    }

    // arctan series: arctan(1/x) = 1/x - 1/(3x^3) + 1/(5x^5) - ...
    public static IEnumerable<BigInteger> ArcTan1OverX(int x, int digits)
    {
        var magnitude = BigInteger.Pow(10, digits);
        var sum = BigInteger.Zero;
        bool sign = true;

        for (int i = 1; true; i += 2)
        {
            var current = magnitude / (BigInteger.Pow(x, i) * i);

            if (sign)
                sum += current;
            else
                sum -= current;

            yield return sum;
            sign = !sign;
        }
    }

    public static void Main(string[] args)
    {
        int digits = 100;
        int iterations = 100;

        if (args.Length > 0 && int.TryParse(args[0], out var parsedDigits) && parsedDigits > 0)
        {
            digits = parsedDigits;
        }

        if (args.Length > 1 && int.TryParse(args[1], out var parsedIterations) && parsedIterations > 0)
        {
            iterations = parsedIterations;
        }

        BigInteger piBigInt = GetPi(digits, iterations);
        string piString = piBigInt.ToString().PadLeft(digits + 1, '0');
        Console.WriteLine(piString.Insert(1, "."));
    }
}
