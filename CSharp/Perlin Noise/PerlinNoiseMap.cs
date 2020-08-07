using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.WebSockets;
using System.Numerics;
using System.Text;

namespace PerlinNoise
{
    public sealed class PerlinNoiseMap
    {
        private int xSize;
        private int ySize;
        private Dictionary<Vector2, Vector2> VectorPairs;

        public static PerlinNoiseMap GenerateNoiseMap(int xSize, int ySize)
        {
            Dictionary<Vector2, Vector2> buffer = new Dictionary<Vector2, Vector2>();
            for(float i = 0f; i <= 1; i+=1.0f/xSize)
            {
                for(float j = 0f; j <= 1; j+=1.0f/ySize)
                {
                    buffer.Add(new Vector2(i, j), new Vector2(RandomExtensions.RandomFloat(), RandomExtensions.RandomFloat()));
                }
            }
            return new PerlinNoiseMap(xSize, ySize, buffer);
        }

        internal void VectorDebug()
        {
            Console.WriteLine("<---------------------DEBUG--------------------->");
            foreach (KeyValuePair<Vector2, Vector2> VectorPair in VectorPairs)
            {
                Console.WriteLine("{0} -> {1}", VectorPair.Key, VectorPair.Value);
            }
            Console.WriteLine("<---------------------DEBUG--------------------->");
        }

        public float Noise(float x, float y)
        {
            if (x < 0f || x > 1f || y < 0f || y > 1f) throw new ArgumentException("The arguments for the noise function must consist of numbers between 0 and 1");
            Vector2 paramPoint = new Vector2(fade(x), fade(y));
            Dictionary<Vector2, float> distances = new Dictionary<Vector2, float>();
            foreach (Vector2 v in VectorPairs.Keys.ToArray())
            {
                distances.Add(v, Vector2.Distance(v, paramPoint));
            }

            Vector2[] closestVecs = distances.OrderBy(X => X.Value).Select(X => X.Key).Take(6).ToArray();
            Vector2 a = closestVecs[0];
            Vector2 b = closestVecs.Where(X => X.X == a.X && X != a).First();
            Vector2 c = closestVecs.Where(X => X.Y == b.Y && X != b).First();
            Vector2 d = closestVecs.Where(X => X.X == c.X && X != c).First();

            closestVecs = new Vector2[]
            {
                a,b,c,d,
            };

            closestVecs = new Vector2[]
            {
                closestVecs.OrderByDescending(X => X.Y).ThenBy(X => X.X).First(),
                closestVecs.OrderByDescending(X => X.Y).ThenByDescending(X => X.X).First(),
                closestVecs.OrderBy(X => X.X).ThenBy(X => X.Y).First(),
                closestVecs.OrderByDescending(X => X.X).ThenBy(X => X.Y).First(),
            };

            Vector2[] gradVecs = new Vector2[]
            {
                VectorPairs[closestVecs[0]],
                VectorPairs[closestVecs[1]],
                VectorPairs[closestVecs[2]],
                VectorPairs[closestVecs[3]],
            };

            float[] dots = new float[4];
            for(int i = 0; i < dots.Length; i++)
            {
                dots[i] = Vector2.Dot(gradVecs[i], closestVecs[i].Distance(paramPoint));
            }

            float AB = dots[0] + x * (dots[1] - dots[0]);
            float CD = dots[2] + x * (dots[3] - dots[2]);

            return AB + y * (CD - AB);
        }

        public float Noise(int X, int Y)
        {
            float Xin = (1f / xSize) * X;
            float Yin = (1f / ySize) * Y;

            if (Xin < 0f) Xin = 0f;
            if (Xin > 1f) Xin = 1f;
            if (Yin < 0f) Yin = 0f;
            if (Yin > 1f) Yin = 1f;
            return Noise(Xin, Yin);
        }

        public static float fade(float t)
        {
            // Fade function as defined by Ken Perlin.  This eases coordinate values
            // so that they will ease towards integral values.  This ends up smoothing
            // the final output.
            return t * t * t * (t * (t * 6 - 15) + 10);         // 6t^5 - 15t^4 + 10t^3
        }

        private PerlinNoiseMap(int xSize, int ySize, Dictionary<Vector2, Vector2> VectorPairs)
        {
            this.xSize = xSize;
            this.ySize = ySize;
            this.VectorPairs = VectorPairs;
        }
    }
}
