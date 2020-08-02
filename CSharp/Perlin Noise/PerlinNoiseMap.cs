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

        // Denne funksjonen genererer et noisemap på med antall vektorer som brukeren spesifiserer
        public static PerlinNoiseMap GenerateNoiseMap(int xSize, int ySize)
        {
            Dictionary<Vector2, Vector2> buffer = new Dictionary<Vector2, Vector2>();
            // Iterere fra 0 til 1 med intervaller på 1/xSize
            for(float i = 0f; i <= 1; i+=1.0f/xSize)
            {
                // Iterere fra 0 til 1 med intervaller på 1/ySize
                for(float j = 0f; j <= 1; j+=1.0f/ySize)
                {
                    // Legge til punktet i buffer dictionary, gi vektoren en tilfeldig pekeretning
                    buffer.Add(new Vector2(i, j), new Vector2(RandomExtensions.RandomFloat(), RandomExtensions.RandomFloat()));
                }
            }
            return new PerlinNoiseMap(xSize, ySize, buffer);
        }

        public float Noise(float x, float y)
        {
            if (x < 0f || x > 1f || y < 0f || y > 1f) throw new ArgumentException("The arguments for the noise function must consist of numbers between 0 and 1");
            // Konvertere parameteret til vektor
            Vector2 paramPoint = new Vector2(fade(x), fade(y));
            // Opprette en dictionary over distanser mellom vektorer og punktet brukeren velger
            Dictionary<Vector2, float> distances = new Dictionary<Vector2, float>();
            foreach (Vector2 v in VectorPairs.Keys.ToArray())
            {
                distances.Add(v, Vector2.Distance(v, paramPoint));
            }

            // Velge de 6 nærmeste punktene
            Vector2[] closestVecs = distances.OrderBy(X => X.Value).Select(X => X.Key).Take(6).ToArray();
            // Det første punktet er alltid det absolutt nærmeste
            Vector2 a = closestVecs[0];
            // Det andre punktet er det nærmeste punktet med samme X verdi som den andre (og som ikke er identisk)
            Vector2 b = closestVecs.Where(X => X.X == a.X && X != a).First();
            // Det tredje punktet er et punkt med samme Y verdi som forrige punkt (og som ikke er identisk)
            Vector2 c = closestVecs.Where(X => X.Y == b.Y && X != b).First();
            // Det tredje punktet er et punkt med samme X verdi som forrige punkt (og som ikke er identisk)
            Vector2 d = closestVecs.Where(X => X.X == c.X && X != c).First();

            // Oversrkive arrayen
            closestVecs = new Vector2[]
            {
                a,b,c,d,
            };

            // Vi må nå sortere lista
            closestVecs = new Vector2[]
            {
                // Punkt A -> punktet med høytest Y verdi, men lavest X
                closestVecs.OrderByDescending(X => X.Y).ThenBy(X => X.X).First(),
                // Punkt B -> punktet med både høyest Y verdi og høyest X verdi
                closestVecs.OrderByDescending(X => X.Y).ThenByDescending(X => X.X).First(),
                // Punkt C punktet med lavest X og lavest Y verdi
                closestVecs.OrderBy(X => X.X).ThenBy(X => X.Y).First(),
                // Punkt D punktet med høyest X verdi og lavest Y verdi
                closestVecs.OrderByDescending(X => X.X).ThenBy(X => X.Y).First(),
            };

            Vector2[] gradVecs = new Vector2[]
            {
                // Finne hvilken retning disse vektorene peker i (lagret i VectorPairs)
                VectorPairs[closestVecs[0]],
                VectorPairs[closestVecs[1]],
                VectorPairs[closestVecs[2]],
                VectorPairs[closestVecs[3]],
            };

            float[] dots = new float[4];
            for(int i = 0; i < dots.Length; i++)
            {
                // Opprette liste over prikkprodukt fra vektorene som peker i tilfeldig retning
                // og distansen mellom vektoren til punktet brukeren putter inn
                dots[i] = Vector2.Dot(gradVecs[i], closestVecs[i].Distance(paramPoint));
            }

            // Lerp
            float AB = dots[0] + x * (dots[1] - dots[0]);
            float CD = dots[2] + x * (dots[3] - dots[2]);



            return fade(AB + y * (CD - AB));
        }

        // Noise funksjon som tar i mot INT (fungerer ikke veldig bra)
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
            // Fade funksjon til Ken Perlin
            return t * t * t * (t * (t * 6 - 15) + 10);
        }

        // Privat konstruktor
        private PerlinNoiseMap(int xSize, int ySize, Dictionary<Vector2, Vector2> VectorPairs)
        {
            this.xSize = xSize;
            this.ySize = ySize;
            this.VectorPairs = VectorPairs;
        }
    }
}
