export namespace cracker {
	
	export class CrackResult {
	    result: string;
	    time: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new CrackResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.result = source["result"];
	        this.time = source["time"];
	        this.error = source["error"];
	    }
	}

}

export namespace database {
	
	export class DecryptResult {
	    save_path: string;
	    wxid: string;
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new DecryptResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.save_path = source["save_path"];
	        this.wxid = source["wxid"];
	        this.err = source["err"];
	    }
	}

}

export namespace extractor {
	
	export class ExtractMapResult {
	    // Go type: orderedmap
	    data?: any;
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new ExtractMapResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], null);
	        this.err = source["err"];
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

}

export namespace winreg {
	
	export class RegResult {
	    // Go type: orderedmap
	    data?: any;
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new RegResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], null);
	        this.err = source["err"];
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

}

