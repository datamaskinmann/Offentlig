using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Management;
using System.Diagnostics;


namespace SystemInfo
{
    static class Management
    {
        // Brukt for å enklere holde styr på System Management Queries
        // Readonly slik at ingen kan endre på verdiene
        // Internal slik at ingen utenfor biblioteket vårt kan aksessere disse verdiene (og dermed endre dem)
        static internal IReadOnlyDictionary<int, string> queries = new Dictionary<int, string>()
        {
            // 1 = info om systemets prosessor
            {1, "select * from Win32_Processor"},
            // 2 = finne kommandolinjeargumenter ved hjelp av prosessid
            {2, "SELECT CommandLine FROM Win32_Process WHERE ProcessId = " },
            // 3 = info om CPU Core 0
            {3, "Win32_Processor.DeviceID='CPU0'"},
            // 4 = info om operativsystem
            {4, "SELECT * FROM Win32_OperatingSystem"},
            // 5 = displayinfo
            {5, "SELECT * FROM Win32_DisplayConfiguration"},
            // 6 = info om videokontroller GPU stats
            {6, "SELECT AdapterRAM from Win32_VideoController"},
        };

        // extractSystemInfoString tar inn som parameter en System management query og henter ut
        // egenskapen definert i selector parameteret
        // dette returneres som en string
        public static string extractSystemInfoString(string query, string selector)
        {
            using (ManagementObjectSearcher mos = new ManagementObjectSearcher(query))
            {
                using (ManagementObjectCollection moc = mos.Get())
                {
                    return moc?.Cast<ManagementBaseObject>()?.SingleOrDefault()?[selector]?.ToString();
                }
            }
        }

        // extractSystemInfoUint tar inn som parameter en System management query og henter ut
        // egenskapen definert i selector parameteret
        // dette returneres som en uint
        public static uint extractSystemInfoUINT(string query, string selector)
        {
            using (var managementObject = new ManagementObject(query))
            {
                var sp = (uint)(managementObject[selector]);

                return sp;
            }
        }

        /// <summary>
        /// Runs a CMD command and reads from echo
        /// </summary>
        /// <param name="args">The command you want to run</param>
        /// <returns>A string array containing each line of the result</returns>
        public static string[] CMDCOMMAND(string args)
        {
            Process p = new Process();
            p.StartInfo = new ProcessStartInfo()
            {
                RedirectStandardOutput = true,
                UseShellExecute = false,
                Arguments = args,
                FileName = "cmd.exe",
            };
            p.Start();
            p.WaitForExit();
            List<string> retArray = new List<string>();
            while (!p.StandardOutput.EndOfStream)
            {
                retArray.Add(p.StandardOutput.ReadLine());
            }
            return retArray.ToArray();
        }
    }
}
