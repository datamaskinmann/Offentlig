using System;
using System.Collections.Generic;
using System.Text;

namespace PerlinNoise
{
    public static class RandomExtensions
    {
        // Hentet fra https://stackoverflow.com/questions/3365337/best-way-to-generate-a-random-float-in-c-sharp
        // Gir oss tall mellom -1 og 1, venter 1 ms
        public static float RandomFloat()
        {
            System.Threading.Thread.Sleep(1);
            return Convert.ToSingle(new Random(DateTime.Now.Millisecond).Next(-100, 100))/ 100;
        }
    }
}
