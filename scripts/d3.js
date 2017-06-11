function setupGraph() {
    // set the dimensions and margins of the graph
    // TODO: Get the device width here?
    var margin = { top: 20, right: 20, bottom: 30, left: 50 },
        width, height;

    width = document.body.clientWidth - margin.left - margin.right - 100;
    height = document.body.scrollHeight - margin.top - margin.bottom - 50;

    // parse the date / time
    // 2017-05-05T12:08:55.745+01:00
    var parseTime = d3.timeParse('%Y-%m-%dT%H:%M:%S.%L%Z');

    // append the svg obgect to the body of the page
    // appends a 'group' element to 'svg'
    // moves the 'group' element to the top left margin
    var svg = d3.select("body").append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform",
            "translate(" + margin.left + "," + margin.top + ")");
    // set the ranges
    var x = d3.scaleTime().range([0, width]);
    var y = d3.scaleLinear().range([height, 0]);

    // define the line
    var valueline = d3.line()
        .x(function(d) { return x(d.time); })
        .y(function(d) { return y(d.count); });

    // Get the data
    d3.json("/data", function(error, data) {
        if (error) throw error;

        var minX = new Date(),
            maxX = 0,
            maxY = 0,
            rangeX = 0,
            maxQueueY = 0;

        series = [];

        data.forEach(function(d1) {
            d1.counts.forEach(function(c) {
                c.time = parseTime(c.time);
                c.count = +c.count;
                if (!c.count) {
                    c.count = 0;
                }
            });

            rangeX = d3.extent(d1.counts, function(d) { return d.time; });

            if (rangeX[1] > maxX) { maxX = rangeX[1]; }
            if (rangeX[0] < minX) { minX = rangeX[0]; }

            maxQueueY = d3.max(d1.counts, function(d) { return d.count; });

            if (maxQueueY > maxY) { maxY = maxQueueY; }

            series.push(d1.counts);
        });

        x.domain([minX, maxX]);
        y.domain([0, maxY]);

        svg.selectAll("path")
            .data(series)
            .enter()
            .append("path")
            .attr("fill", "none")
            .attr("stroke", "steelblue")
            .attr("stroke-linejoin", "round")
            .attr("stroke-linecap", "round")
            .attr("stroke-width", 1.5)
            .attr("d", valueline);

        // Add the valueline path.

        // Add the X Axis
        svg.append("g")
            .attr("transform", "translate(0," + height + ")")
            .call(d3.axisBottom(x));

        // Add the Y Axis
        svg.append("g")
            .call(d3.axisLeft(y));
    });
}

function getValue(text) {
    if (!text) {
        return 0;
    }
    var k = text.indexOf("k");

    // TODO: other units?
    if (k == -1) {
        return +text;
    }

    return text.slice(0, k) * 1000;
}