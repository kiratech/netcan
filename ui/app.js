var createGraph = function (elements) {
    return cytoscape({
        container: document.getElementById('cy'),

        boxSelectionEnabled: false,
        autounselectify: false,

        style: cytoscape.stylesheet()
            .selector('node')
            .css({
                'content': 'data(id)'
            })
            .selector('edge')
            .css({
                'target-arrow-shape': 'triangle',
                'width': 4,
                'line-color': '#ddd',
                'target-arrow-color': '#ddd',
                'curve-style': 'bezier'
            })
            .selector('.highlighted')
            .css({
                'background-color': '#61bffc',
                'line-color': '#61bffc',
                'target-arrow-color': '#61bffc',
                'transition-property': 'background-color, line-color, target-arrow-color',
                'transition-duration': '0.5s'
            }),

        elements: elements,
        layout: {
            name: 'breadthfirst',
            directed: true,
            padding: 10
        }
    });
};


var main = function () {
    const httpBaseUrl = window.location.origin
    const websocketBaseUrl = "ws://" + httpBaseUrl.replace(/.*?:\/\//g, "");
    var ws = new WebSocket(websocketBaseUrl + "/ws");
    ws.onmessage = function (evt) {
        var result = JSON.parse(evt.data);
        //TODO: find a better way to update than recreating the graph
        createGraph(result.elements);
    };
};

window.onload = main;

