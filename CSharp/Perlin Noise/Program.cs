using System;
using System.Diagnostics;
using System.Reflection.Metadata.Ecma335;
using System.Runtime.Serialization;

namespace PerlinNoise
{
    class Program
    {
        static PerlinNoiseMap m;
        static void Main(string[] args)
        {
            Console.SetBufferSize(Console.LargestWindowWidth, Console.LargestWindowHeight);
            while (true)
            {
                m = PerlinNoiseMap.GenerateNoiseMap(1, 1);
                for (int i = 0; i < Console.LargestWindowWidth - 1; i++)
                {
                    for (int j = 0; j < Console.LargestWindowHeight - 1; j++)
                    {
                        Console.SetCursorPosition(i, j);
                        WriteColour("■", GetPixelColour(i, j));
                    }
                }
                switch (Console.ReadKey().Key)
                {
                    case ConsoleKey.Enter: continue;
                    case ConsoleKey.Backspace:
                        m.VectorDebug();
                        Console.ReadLine();
                        break;
                    case ConsoleKey.DownArrow:
                        Debug.WriteLine(Console.CursorLeft + " " + Console.CursorTop);
                        break;
                    case ConsoleKey.Escape: break;
                }
            }
        }

        static void WriteColour(string Text, ConsoleColor Color)
        {
            Console.ForegroundColor = Color;
            Console.Write(Text);
            Console.ResetColor();
        }

        static ConsoleColor GetPixelColour(int left, int top)
        {
            float result = 0f;
            if(left == 0 && top == 0)
            {
                result = m.Noise(0f, 0f);
            }
            else if(left == 0)
            {
                result = m.Noise(0f, (1f /Console.LargestWindowHeight)*top);
            }
            else if(top == 0)
            {
                result = m.Noise((1f /Console.LargestWindowWidth)*left, 0f);
            }
            else
            {
                result = m.Noise((1f / Console.LargestWindowWidth)*left, (1f / Console.LargestWindowHeight)*top);
            }
            result = Math.Abs(result);
            if (result > 1f) result = 1f;
            return result < 0.15 ? ConsoleColor.Blue
                : result > 0.15 && result < 0.30 ? ConsoleColor.Yellow
                : result > 0.30 && result < 0.75 ? ConsoleColor.DarkGreen
                : result > 0.75 && result < 0.75125 ? ConsoleColor.Cyan
                : ConsoleColor.DarkGray;
        }
    }
}
