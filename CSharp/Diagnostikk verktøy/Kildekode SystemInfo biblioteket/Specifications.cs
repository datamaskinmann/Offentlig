using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Management;
using System.Diagnostics;
using System.Reflection;

// Forfatter: Sven Sørensen
// For MP03
// Opprettet: 09 Mars 2020
namespace SystemInfo
{
    /// <summary>
    /// This class contains specifications associated with this computer
    /// </summary>
    public static class Specifications
    {
        // ProcessorType representerer typen prosessor systemet ditt har  
        /// <summary>
        /// Gets the processor-model of this computer
        /// </summary>
        public static string ProcessorModel { get { return Management.extractSystemInfoString(Management.queries[1], "name"); } }
        // Oppgave 2 G
        // ProsessorCores representerer antall cores prosessoren har
        /// <summary>
        /// Gets the amount of cores in this computer's processor
        /// </summary>
        public static ushort ProcessorCores
        {
            get
            {
                return Convert.ToUInt16(Management.extractSystemInfoString(Management.queries[1], "NumberOfCores"));
            }
        }

        // ProcessorArchitecture representerer hvilken arkitektur prosessoren din har
        // samt hvor mange biter den har
        /// <summary>
        /// Gets the processor architecture of this computer
        /// </summary>
        public static string ProcessorArchitecture
        {
            get
            {
                // Finne prosessorarkitekturen
                return Management.CMDCOMMAND("/C echo %PROCESSOR_ARCHITECTURE%")[0];
            }
        }

        // Oppgave 2 F
        // ProcessorMaxClockFrequeny representerer hvor høy klokkehastighet prosessoren din har
        /// <summary>
        /// Gets the max clock frequency of this computer's processor in megaherz
        /// </summary>
        public static uint ProcessorMaxClockFrequency { get { return Management.extractSystemInfoUINT(Management.queries[3], "MaxClockSpeed"); } }
        // L1CacheSize representerer størrelsen på cache L1 i prosessoren din
        /// <summary>
        /// Gets the L1 cache size of this computer's processor
        /// </summary>
        public static uint L1CacheSize { get { return cacheSize(3); } }
        // L2CacheSize representerer størrelsen på cahce L2 i prosessoren din
        /// <summary>
        /// Gets the L2 cache size of this computer's processor
        /// </summary>
        public static uint L2CacheSize { get { return cacheSize(4); } }
        // L3CacheSize representerer størrelsen på cahce L3 i prosessoren din
        /// <summary>
        /// Gets the L3 cache size of this computer's processor
        /// </summary>
        public static uint L3CacheSize { get { return cacheSize(5); } }

        // VisibleMemorySize representerer totalt synlig minne (RAM) på din datamaskin
        /// <summary>
        /// Gets the total visible memory of this computer in bytes
        /// </summary>
        public static ulong TotalMemorySize
        {
            get
            {
                string[] rawMem = Management.CMDCOMMAND("/C WMIC MEMORYCHIP get CAPACITY");
                ulong total = 0;
                for(int i = 1; i < rawMem.Length; i++)
                {
                    ulong num;
                    ulong.TryParse(rawMem[i], out num);
                    total += num;
                }
                return total;
            }
        }

        // GPUInfo representerer  en beskrivelse av datamaskinens grafikk-kort
        /// <summary>
        /// Gets the GPU Model of this computer
        /// </summary>
        public static string GPUModel { get { return Management.extractSystemInfoString(Management.queries[5], "Description"); } }

        // GPUVram representerer datamaskinens VRAM i antall bytes
        /// <summary>
        /// Gets the amount of dedicated vram of this computer's GPU in bytes
        /// </summary>
        public static ulong GPUVRAM
        {
            get
            {
                return Convert.ToUInt64(Management.extractSystemInfoString(Management.queries[6], "AdapterRam"));
            }
        }

        // NB: Koden for CacheSize metoden er
        // inspirert fra Stackoverflow på
        // https://stackoverflow.com/questions/6995787/how-to-determine-cpu-cache-size-in-net

        // cacheSize returnerer størrelsen på cache oppgitt i parameteret layer
        // i antall KB
        private static uint cacheSize(int layer)
        {
            // Av en eller annen grunn representeres cache L1 som 3, L2 som 4 og L3 som 5
            // Dermed kaster vi en ArgumentException dersom layer parameteret ikke er innenfor
            // disse begrensningene
            if (layer < 3 || layer > 5) throw new ArgumentException("The value '" + layer + "' is invalid. Use 3 (L1), 4 (L2) or 5 (L3)");
            using (ManagementClass mc = new ManagementClass("Win32_CacheMemory"))
            {
                using (ManagementObjectCollection moc = mc.GetInstances())
                {
                    // LINQ Query som finner størrelsen på den ønskede cachen
                    // oppgitt i parameteret layer
                    return Convert.ToUInt32(moc.Cast<ManagementObject>()
                        .Where(X => (ushort)
                        X.Properties["Level"]
                        .Value == (ushort)layer)
                        .FirstOrDefault().
                        Properties["MaxCacheSize"].Value);
                }
            }
        }
    }
}
