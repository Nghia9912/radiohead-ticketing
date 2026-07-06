import http from 'k6/http';
import { check, sleep } from 'k6';

// Simulate 5000 concurrent users joining the waiting room over 30 seconds, holding for 1 minute
export const options = {
    stages: [
        { duration: '10s', target: 50000 }, // Ramp up to 50000 users
        { duration: '30s', target: 50000 }, // Stay at 50000 users
        { duration: '10s', target: 0 },    // Ramp down to 0
    ],
};

export default function () {
    // In a real scenario, this would be hitting /api/v1/queue/enqueue
    // For now, we hit the /metrics endpoint to simulate load on the HTTP server
    const res = http.get('http://host.docker.internal:8080/metrics');
    
    check(res, {
        'status is 200': (r) => r.status === 200,
    });
    
    sleep(1);
}
