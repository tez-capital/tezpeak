import type { ApplicationServices } from "@src/common/types/status";

export function extractContinualServiceInfo(services?: ApplicationServices) {
    if (!services) return undefined;

    for (const serviceName in services) {
        if (serviceName.endsWith("-continual")) {
            return services[serviceName];
        }
    }
    return undefined
}