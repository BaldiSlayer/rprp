#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <omp.h>

void simple_iteration(double *A, double *b, double *x0, double tau, double epsilon, int n, int max_iter) {
    double *x = (double *)malloc(n * sizeof(double));
    double *Ax = (double *)malloc(n * sizeof(double));
    double *r = (double *)malloc(n * sizeof(double));

    for (int i = 0; i < n; ++i) {
        x[i] = x0[i];
    }

    for (int iter = 0; iter < max_iter; ++iter) {
        #pragma omp parallel for
        for (int i = 0; i < n; ++i) {
            Ax[i] = 0;

            for (int j = 0; j < n; ++j) {
                Ax[i] += A[i * n + j] * x[j];
            }
        }

        double norm_r = 0.0;
        double norm_b = 0.0;

        #pragma omp parallel for reduction(+:norm_r, norm_b)
        for (int i = 0; i < n; ++i) {
            r[i] = Ax[i] - b[i];

            norm_r += r[i] * r[i];
            norm_b += b[i] * b[i];
        }

        if (sqrt(norm_r) < epsilon * sqrt(norm_b)) {
            break;
        }

        #pragma omp parallel for
        for (int i = 0; i < n; ++i) {
            x[i] = x[i] - tau * r[i];
        }
    }

    free(Ax);
    free(r);
    free(x);
}

signed main(int argc, char *argv[]) {
    int n = 8000;
    double tau = 0.1 / n;
    double epsilon = 1e-5;

    int max_iter = 1500;

    double *A = (double *)malloc(n * n * sizeof(double));
    double *b = (double *)malloc(n * sizeof(double));
    double *x0 = (double *)malloc(n * sizeof(double));

    #pragma omp parallel for
        for (int i = 0; i < n; ++i) {
            b[i] = n + 1;

            for (int j = 0; j < n; ++j) {
                A[i * n + j] =  1.0 + (i == j);
            }
        }

    #pragma omp parallel for
        for (int i = 0; i < n; i++) {
            x0[i] = 0.0;
        }

    double start_time = omp_get_wtime();

    simple_iteration(A, b, x0, tau, epsilon, n, max_iter);

    double end_time = omp_get_wtime();

    printf("Время выполнения: %.6f секунд\nИспользовано потоков: %d\n", end_time - start_time, omp_get_max_threads());

    free(A);
    free(b);
    free(x0);

    return 0;
}
