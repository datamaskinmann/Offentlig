using System;
using System.Collections.Generic;
using System.Text;
using System.Threading;

namespace PerlinNoise
{
    public static class RandomExtensions
    {
        // Hentet fra https://stackoverflow.com/questions/3365337/best-way-to-generate-a-random-float-in-c-sharp
        static int last = 0;
        public static float RandomPositiveFloat()
        {
            Thread.Sleep(1);
            int rnd = new Random(DateTime.Now.Millisecond+last).Next(0, 10000000);
            last = rnd;
            if (rnd == 0) return 0f;
            return Convert.ToSingle(rnd / 10000000f);
        }
        public static float RandomNegativeFloat()
        {
            Thread.Sleep(1);
            int rnd = new Random(DateTime.Now.Millisecond+last).Next(-10000000, 0);
            last = rnd;
            if (rnd == 0) return 0f;
            return Convert.ToSingle(rnd / 10000000f);
        }

        public static float RandomFloat()
        {
            Thread.Sleep(1);
            int rnd = new Random(DateTime.Now.Millisecond + last).Next(-10000000, 10000000);
            if (rnd == 0) return 0f;
            return Convert.ToSingle(rnd / 10000000f);
        }
    }
}
