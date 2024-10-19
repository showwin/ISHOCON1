interface TeamScore {
    team: string;
    score: number;
    timestamp: string;
}

const apiUrl = "<<API_GATEWAY_DOMAIN_URL>>teams";

// Function to fetch the latest data
async function fetchData() {
    try {
        const response = await fetch(apiUrl);
        const data = await response.json();
        return data;
    } catch (error) {
        console.error("Error fetching data:", error);
        return [];
    }
}

// Define color scale globally to use across both graphs
const colorScale = d3.scaleOrdinal(d3.schemeTableau10);

function renderTimeline(data: TeamScore[]) {
    const margin = { top: 40, right: 30, bottom: 50, left: 80 };  // Adjust left margin for larger Y-axis values
    const width = document.body.clientWidth * 0.8 - margin.left - margin.right;
    const height = 400 - margin.top - margin.bottom;

    // Create a title for the graph
    d3.select("#timeline-chart")
        .append("h2")
        .style("text-align", "center")
        .text("Timeline");

    const svg = d3.select("#timeline-chart")
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform", `translate(${margin.left},${margin.top})`);

    const parseDate = d3.timeParse("%Y-%m-%dT%H:%M:%S%Z");
    data.forEach(d => {
        d.timestamp = parseDate(d.timestamp)?.toString() ?? "";
    });

    data = data.filter(d => !isNaN(d.score) && d.timestamp);
    const dates = data.map(d => new Date(d.timestamp));

    const x = d3.scaleTime()
        .domain([d3.min(dates)!, d3.max(dates) || new Date()])
        .range([0, width]);

    const maxYValue = d3.max(data, d => d.score) || 0;
    const y = d3.scaleLinear()
        .domain([0, maxYValue * 1.1])  // Add 10% padding above the maximum value
        .range([height, 0]);

    // X-axis with default time format (HH:MM)
    svg.append("g")
        .attr("transform", `translate(0,${height})`)
        .call(d3.axisBottom(x).tickSize(-height).tickPadding(10));

    // Y-axis with full numbers (no abbreviation)
    svg.append("g")
        .call(d3.axisLeft(y).tickSize(-width).tickPadding(10));

    // Lighter grid lines
    svg.selectAll(".tick line")
        .attr("stroke", "#ddd")
        .attr("stroke-opacity", 0.3);

    // Line path for each team
    const line = d3.line<TeamScore>()
        .x(d => x(new Date(d.timestamp)))
        .y(d => y(d.score))
        .curve(d3.curveLinear);

    const teams = d3.group(data, d => d.team);

    teams.forEach((teamData, teamName) => {
        // Draw the line for each team
        svg.append("path")
            .datum(teamData)
            .attr("fill", "none")
            .attr("stroke", colorScale(teamName)!)
            .attr("stroke-width", 2)
            .attr("d", line);

        // Find the last data point for the team
        const lastDataPoint = teamData.reduce((latest, current) =>
            new Date(current.timestamp) > new Date(latest.timestamp) ? current : latest
        );

        // Draw a horizontal line extending from the last data point to the right edge of the graph
        svg.append("line")
            .attr("x1", x(new Date(lastDataPoint.timestamp)))  // Start at the last data point
            .attr("x2", width)  // Extend to the right edge of the graph
            .attr("y1", y(lastDataPoint.score))  // Y position based on the last score
            .attr("y2", y(lastDataPoint.score))  // Keep Y constant to form a horizontal line
            .attr("stroke", colorScale(teamName)!)  // Use the team's color for the line
            .attr("stroke-width", 2)
            .attr("stroke-opacity", 0.3)  // Lighter line using reduced opacity
            .attr("stroke-dasharray", "4,4");  // Dashed line for differentiation
    });

    // Tooltip container
    const tooltip = d3.select("body").append("div").attr("class", "tooltip-modern").style("display", "none");

    // Circles at data points and mouseover event for tooltip (showing only HH:MM:SS)
    const timeTooltipFormat = d3.timeFormat("%H:%M:%S");
    svg.selectAll("dot")
        .data(data)
        .enter()
        .append("circle")
        .attr("cx", d => x(new Date(d.timestamp)))
        .attr("cy", d => y(d.score))
        .attr("r", 5)
        .attr("fill", d => colorScale(d.team)!)
        .on("mouseover", function (event, d) {
            tooltip.style("display", "block")
                .html(`Team: ${d.team}<br>Score: ${d.score}<br>Time: ${timeTooltipFormat(new Date(d.timestamp))}`)
                .style("left", (event.pageX + 10) + "px")
                .style("top", (event.pageY - 20) + "px");
        })
        .on("mouseout", function () {
            tooltip.style("display", "none");
        });
}

function renderBarChart(data: TeamScore[]) {
    const margin = { top: 40, right: 30, bottom: 50, left: 100 };
    const width = document.body.clientWidth * 0.8 - margin.left - margin.right;
    const height = 400 - margin.top - margin.bottom;

    // Sort by timestamp first, then filter out old scores, keeping only the latest for each team
    const latestScores = Array.from(d3.group(data, d => d.team).values()).map(teamScores => {
        return teamScores.reduce((latest, current) => {
            return new Date(current.timestamp) > new Date(latest.timestamp) ? current : latest;
        });
    });

    // Sort the latest scores by score in descending order
    latestScores.sort((a, b) => b.score - a.score);

    // Create a title for the graph
    d3.select("#bar-chart")
        .append("h2")
        .style("text-align", "center")
        .text("Latest Score");

    const svg = d3.select("#bar-chart")
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform", `translate(${margin.left},${margin.top})`);

    const y = d3.scaleBand()
        .domain(latestScores.map(d => d.team))
        .range([0, height])
        .padding(0.1);

    const maxYValue = d3.max(latestScores, d => d.score) || 0;
    const x = d3.scaleLinear()
        .domain([0, maxYValue * 1.1])  // Add 10% padding above the maximum value
        .range([0, width]);

    // X-axis
    svg.append("g")
        .attr("transform", `translate(0,${height})`)
        .call(d3.axisBottom(x).tickSize(-height).tickPadding(10));

    // Y-axis
    svg.append("g")
        .call(d3.axisLeft(y).tickSize(-width).tickPadding(10));

    svg.selectAll(".tick line").attr("stroke", "#ddd").attr("stroke-opacity", 0.3);

    const tooltip = d3.select("body").append("div").attr("class", "tooltip-modern").style("display", "none");

    svg.selectAll(".bar")
        .data(latestScores)
        .enter()
        .append("rect")
        .attr("x", 0)
        .attr("y", d => y(d.team)!)
        .attr("width", d => x(d.score))
        .attr("height", y.bandwidth())
        .attr("fill", d => colorScale(d.team)!)
        .on("mouseover", function (event, d) {
            tooltip.style("display", "block")
                .html(`Team: ${d.team}<br>Score: ${d.score}`)
                .style("left", (event.pageX + 10) + "px")
                .style("top", (event.pageY - 20) + "px");
        })
        .on("mouseout", function () {
            tooltip.style("display", "none");
        });
}

// Function to update the graph
function updateGraph() {
    fetchData().then(data => {
        // Clear existing graphs
        d3.select("#timeline-chart").selectAll("*").remove();
        d3.select("#bar-chart").selectAll("*").remove();

        // Re-render the graphs with the latest data
        renderTimeline(data);
        renderBarChart(data);
    });
}

// Initial graph rendering
updateGraph();

// Fetch and refresh the graph every 3 minutes (180,000 ms)
setInterval(updateGraph, 180000);
