export namespace main {
	
	export class Message {
	    role: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	    }
	}
	export class Model {
	    name: string;
	    model: string;
	    modified_at: string;
	    size: number;
	    digest: string;
	    details: Record<string, any>;
	    is_running: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Model(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.model = source["model"];
	        this.modified_at = source["modified_at"];
	        this.size = source["size"];
	        this.digest = source["digest"];
	        this.details = source["details"];
	        this.is_running = source["is_running"];
	    }
	}
	export class OllamaServerConfig {
	    id: string;
	    name: string;
	    base_url: string;
	    api_key: string;
	    is_active: boolean;
	    test_status: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new OllamaServerConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.base_url = source["base_url"];
	        this.api_key = source["api_key"];
	        this.is_active = source["is_active"];
	        this.test_status = source["test_status"];
	        this.type = source["type"];
	    }
	}

}

