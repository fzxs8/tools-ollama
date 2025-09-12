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
	export class Conversation {
	    id: string;
	    title: string;
	    messages: Message[];
	    modelName: string;
	    systemPrompt: string;
	    modelParams: string;
	    timestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new Conversation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.messages = this.convertValues(source["messages"], Message);
	        this.modelName = source["modelName"];
	        this.systemPrompt = source["systemPrompt"];
	        this.modelParams = source["modelParams"];
	        this.timestamp = source["timestamp"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
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
	    }
	}
	export class Prompt {
	    id: string;
	    name: string;
	    content: string;
	    description: string;
	    createdAt: number;
	    updatedAt: number;
	    models: string[];
	    version: number;
	    tags: string[];
	    createdBy: string;
	
	    static createFrom(source: any = {}) {
	        return new Prompt(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.content = source["content"];
	        this.description = source["description"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.models = source["models"];
	        this.version = source["version"];
	        this.tags = source["tags"];
	        this.createdBy = source["createdBy"];
	    }
	}

}

