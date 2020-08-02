using System;

namespace PerlinNoise
{
    class Program
    {
        static void Main(string[] args)
        {
            PerlinNoiseMap m = PerlinNoiseMap.GenerateNoiseMap(1, 1);
            Console.WriteLine("Done");
            for (int x = 0; x < 119;)
            {
                for (int y = 0; y < 30;)
                {
                    for (float i = 0f; i < 1f; i += 1 / 119f)
                    {
                        for (float j = 0f; j < 1f; j += 1 / 30f)
                        {
                            float b = m.Noise(i, j);
                            Console.SetCursorPosition(x, y);
                            if (b < 0.3) WriteColour("■", ConsoleColor.Blue);
                            if (b > 0.3 && b < 0.35) WriteColour("■", ConsoleColor.Yellow);
                            if (b > 0.35 && b < 0.75) WriteColour("■", ConsoleColor.DarkGreen);
                            else WriteColour("■", ConsoleColor.DarkGray);
                            y++;
                        }
                        x++;
                        y = 0;
                    }
                    break;
                }
            }
            Console.ReadLine();
        }

        static void WriteColour(string Text, ConsoleColor Color)
        {
            Console.ForegroundColor = Color;
            Console.Write(Text);
            Console.ResetColor();
        }
    }
}
