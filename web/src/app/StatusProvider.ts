export type StatusProviderStatus = "connected" | "disconnected" | "reconnecting" ;

export class StatusProvider {
	private url: string;
	private eventSource: EventSource | null = null;
	private retryCount = 0;
	private retryDelay: number; // Delay in milliseconds

	public onmessage: (event: MessageEvent) => void = () => { };
	public onstatuschange: (status: StatusProviderStatus) => void = () => { };

	constructor(url: string, retryDelay = 5000) {
		this.url = url;
		this.retryDelay = retryDelay;
		this.connect();
	}

	private connect() {
		this.eventSource = new EventSource(this.url);

		this.eventSource.onopen = () => {
			this.retryCount = 0; // Reset retry count on successful connection
			this.onstatuschange('connected');
		};

		this.eventSource.onmessage = (event) => this.onmessage(event);

		this.eventSource.onerror = () => {
			this.onstatuschange("disconnected");
			this.eventSource?.close(); // Close the existing connection

			setTimeout(() => {
				this.retryCount++;
				this.onstatuschange("reconnecting");
				this.connect(); // Attempt to reconnect
			}, this.retryDelay);
		};
	}

	public close() {
		this.eventSource?.close();
	}
}