using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Diagnostics;

namespace SystemInfo
{
    /// <summary>
    /// This class contains diagnostics associated with this computer
    /// </summary>
    public static class Diagnostics
    {
        /// <summary>
        /// Gets the amount of processes with the 'running' state on this computer
        /// </summary>
        public static Process[] RunningProcesses
        {
            get
            {
                string[] PIDS = Management.CMDCOMMAND("/C for /f \"tokens=1,2\" %a in " +
                      "('Tasklist /fi \"STATUS eq running\" /nh') do @echo %b");
                List<Process> retArray = new List<Process>();
                foreach(string s in PIDS)
                {
                    try // Try fordi en prosess kan avslutte før eller mens vi itererer
                    // og da ikke gi utslag på Process.GetProcessById
                    {
                        retArray.Add(Process.GetProcessById(Convert.ToInt32(s)));
                    }
                    catch
                    {
                        continue;
                    }
                }
                return retArray.ToArray();
            }
        }

        /// <summary>
        /// Gets the processes on this system
        /// </summary>
        public static Process[] Processes
        {
            get
            {
                return Process.GetProcesses();
            }
        }
        // MostMemoryIntensive returnerer den prosessen som er mest minneintensiv
        /// <summary>
        /// Gets the most memory-intensive process on this computer
        /// </summary>
        public static Process MostMemoryIntensive
        {
            get
            {
                // LINQ Query som sorterer alle prosessene etter
                // synkende rekkefølge i WorkingSet64
                // Returnerer den første av disse (prosessen som bruker mest ram)
                return Process.GetProcesses()
                    ?.Where(X => X.BasePriority != 0) // Program har ikke tilgang til prosesser med prioritet 0
                    ?.OrderByDescending(Y => Y.WorkingSet64)
                    ?.First();
            }
        }

        // Metoden tar inn Prosess P som parameter og returnerer kommandolinjen for prosessen
        // CommandLine tar inn Process P som parameter og returnerer kommandolinjeargumentene for
        // denne prosessen
        /// <summary>
        /// Gets the commandline-arguments associated with a process
        /// </summary>
        /// <param name="P">The process you want to check the commandline-arguments associated with</param>
        /// <returns>A string which represents the commandline-arguments associated with process P</returns>
        public static string CommandLine(Process P)
        {
            return Management.extractSystemInfoString(Management.queries[2] + P.Id, "CommandLine");
        }

        // CPUUsage returnerer hvor mange % av prosessorkraften
        // prosess P bruker
        /// <summary>
        /// Gets the CPUUsage of a process
        /// </summary>
        /// <param name="P">The process you want to check the CPU usage for</param>
        /// <returns>A double which represents in %, the CPU usage associated with process P</returns>
        public static double CPUUsage(Process P)
        {
            // Where clause som velger alle prosessene og
            // kumulativt adderer totale millisekunder
            // disse prosessene har brukt av prosessoren
            return (P.TotalProcessorTime.TotalMilliseconds /
                Process.GetProcesses()
                .Where(X => X.BasePriority != 0) // Program har ikke tilgang til prosesser med BasePriority 0
                .Select(Y => Y.TotalProcessorTime.TotalMilliseconds)
                .Sum()) * 100;
        }
    }
}
