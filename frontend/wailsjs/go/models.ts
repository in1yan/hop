export namespace main {
	
	export class Window {
	    Handle: any;
	    Title: string;
	
	    static createFrom(source: any = {}) {
	        return new Window(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Handle = source["Handle"];
	        this.Title = source["Title"];
	    }
	}

}

