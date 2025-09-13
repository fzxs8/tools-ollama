export namespace types {
	
	export class FormDataItem {
	    key: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new FormDataItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
	    }
	}
	export class RequestBody {
	    type: string;
	    rawContent?: string;
	    rawContentType?: string;
	    formData?: FormDataItem[];
	
	    static createFrom(source: any = {}) {
	        return new RequestBody(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.rawContent = source["rawContent"];
	        this.rawContentType = source["rawContentType"];
	        this.formData = this.convertValues(source["formData"], FormDataItem);
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
	export class RequestHeader {
	    key: string;
	    value: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RequestHeader(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
	        this.enabled = source["enabled"];
	    }
	}
	export class QueryParam {
	    key: string;
	    value: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new QueryParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
	        this.enabled = source["enabled"];
	    }
	}
	export class ApiRequest {
	    method: string;
	    selectedServerId: string;
	    path: string;
	    queryParams: QueryParam[];
	    headers: RequestHeader[];
	    body: RequestBody;
	
	    static createFrom(source: any = {}) {
	        return new ApiRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.method = source["method"];
	        this.selectedServerId = source["selectedServerId"];
	        this.path = source["path"];
	        this.queryParams = this.convertValues(source["queryParams"], QueryParam);
	        this.headers = this.convertValues(source["headers"], RequestHeader);
	        this.body = this.convertValues(source["body"], RequestBody);
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
	export class ApiResponse {
	    statusCode: number;
	    statusText: string;
	    headers: RequestHeader[];
	    body: string;
	    requestDurationMs: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ApiResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.statusCode = source["statusCode"];
	        this.statusText = source["statusText"];
	        this.headers = this.convertValues(source["headers"], RequestHeader);
	        this.body = source["body"];
	        this.requestDurationMs = source["requestDurationMs"];
	        this.error = source["error"];
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
	export class Message {
	    role: string;
	    content: string;
	    timestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	        this.timestamp = source["timestamp"];
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
	    modifiedAt: string;
	    size: number;
	    digest: string;
	    details: Record<string, any>;
	    isRunning: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Model(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.model = source["model"];
	        this.modifiedAt = source["modifiedAt"];
	        this.size = source["size"];
	        this.digest = source["digest"];
	        this.details = source["details"];
	        this.isRunning = source["isRunning"];
	    }
	}
	export class OllamaServerConfig {
	    id: string;
	    name: string;
	    baseUrl: string;
	    apiKey: string;
	    isActive: boolean;
	    testStatus: string;
	
	    static createFrom(source: any = {}) {
	        return new OllamaServerConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.baseUrl = source["baseUrl"];
	        this.apiKey = source["apiKey"];
	        this.isActive = source["isActive"];
	        this.testStatus = source["testStatus"];
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

