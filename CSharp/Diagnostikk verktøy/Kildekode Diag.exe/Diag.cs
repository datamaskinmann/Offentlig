using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Management;
using System.Diagnostics;
using System.Reflection;
using SystemInfo;

namespace Diag
{
    class Diag
    {
        static void Main(string[] args)
        {
            // Kildekoden til Diagnostics og Specifications finner du under "Kildekode SystemInfo biblioteket"

            Console.WriteLine("Assembly versjon: ({0})\n",
                Assembly.GetEntryAssembly().GetName().Version);

            Console.WriteLine("Prosesser på denne PC'en: ({0})\n", Diagnostics.Processes.Length);
            Console.WriteLine("Prosessser på denne PC'en som er i 'running' tilstand: ({0})\n", Diagnostics.RunningProcesses.Length);

            Process mostMemoryIntensive = Diagnostics.MostMemoryIntensive;
            StringBuilder sb = new StringBuilder();
            sb.Append(String.Format("Den mest minneintensive prosessen på din datamaskin heter "
                + "({0}) og har prossess ID ({1}).\nProsessen bruker ~({2} MiB) minne og ",
                mostMemoryIntensive.ProcessName,
                mostMemoryIntensive.Id,
                mostMemoryIntensive.WorkingSet64
                / Convert.ToInt32(Math.Pow(2, 20)))); // Bytes -> Mib = X / 2^20
            double CPUUsage = Diagnostics.CPUUsage(mostMemoryIntensive);
            sb.Append(
                CPUUsage < 1 ? "< 1%" :
                String.Format("(~{0}%)",
                Convert.ToUInt16(CPUUsage)));

            sb.Append(String.Format(" prosessorkraft.\nKommandolinjeargumentene for denne prosessen er: ({0})\n",
                Diagnostics.CommandLine(mostMemoryIntensive)));

            Console.WriteLine(sb.ToString());

            Console.WriteLine("Din datamaskin har følgende prosessor: ({0}).\n"
                + "Denne prosessoren har følgende arkitektur: ({1}) med en max klokkefrekvens på ({2} MHz)\n",
                Specifications.ProcessorModel,
                Specifications.ProcessorArchitecture,
                Specifications.ProcessorMaxClockFrequency);
            Console.WriteLine("Størrelsen på cache L1, L2 og L3 er påfølgende verdier ({0} KB, {1} KB og {2} KB)\n",
                Specifications.L1CacheSize,
                Specifications.L2CacheSize,
                Specifications.L3CacheSize);
            Console.WriteLine("Prosessoren din har ({0}) cores\n",
                Specifications.ProcessorCores);
            Console.WriteLine("Datamaskinen din har ({0} MB) synlig primært minne (RAM)\n",
                // Konvertering fra bytes til MB ->
                // X/(2^20)
                Convert.ToUInt64(Specifications.TotalMemorySize
                / Math.Pow(2, 20)));
            Console.WriteLine("Datamaskinene din har følgende grafikk-kort: ({0}) med ({1} MB) dedikert VRAM\n",
                Specifications.GPUModel,
                Specifications.GPUVRAM /
                // Konvertering fra Bytes til MB ->
                // X/(2^20)
                Convert.ToUInt32(Math.Pow(2, 20)));
            Console.WriteLine("Trykk enter for å avslutte...");
            Console.ReadLine(); // Hindre programmet fra å umiddelbart avslutte
        }
    }
}