#include "soxr.h"
#include <vector>

#if _MSC_VER // If we are compiling with the Microsoft compiler
    #define EXPORT_API __declspec(dllexport)
#else // For macOS/Linux (Clang/GCC)
    #define EXPORT_API
#endif

// Define a simple structure to hold the resampler state
extern "C" {
    // Factory function to create a resampler instance
    EXPORT_API void* CreateSoxr(double inRate, double outRate, int channels) {
        soxr_error_t error;
        // Using 'soxr_create' from the library
        soxr_t s = soxr_create(inRate, outRate, channels, &error, NULL, NULL, NULL);
        return (void*)s;
    }

    // The main processing function
    EXPORT_API int ProcessSoxr(void* resampler, float* input, int inLen, float* output, int outLen) {
        soxr_t s = (soxr_t)resampler;
        size_t idone, odone;

        // Process the buffers
        // input: pointer to Unity's raw float array
        // output: pointer to the destination buffer
        soxr_error_t error = soxr_process(s, input, (size_t)inLen, &idone, output, (size_t)outLen, &odone);

        return (int)odone; // Return how many samples were actually generated
    }

    EXPORT_API void DestroySoxr(void* resampler) {
        soxr_delete((soxr_t)resampler);
    }

    EXPORT_API int OneshotSoxr(
        double inRate,
        double outRate,
        float* input,
        int inLen,
        float* output,
        int outLen,
        int channels) {
        size_t odone;

        soxr_error_t error = soxr_oneshot(
            inRate, outRate, channels,       /* Rates and # of chans. */
            input, inLen, NULL,       /* Input. */
            output, outLen, &odone,     /* Output. */
            NULL, NULL, NULL);

        return (int)odone;
    }
}