using System;
using System.Collections.Generic;
using System.Numerics;
using System.Text;

namespace PerlinNoise
{
    public static class Vector2Extensions
    {
        public static Vector2 Distance(this Vector2 a, Vector2 b)
        {
            return new Vector2(a.X - b.X, a.Y - b.Y);
        }
    }
}
