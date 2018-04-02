import React, { Component } from "react";
import moment from "moment";

class MareyDiagram extends Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      margins: [75, 20, 100, 140]
    };
  }

  componentDidMount() {
    this.update();
    window.setInterval(this.update.bind(this), 10000);
  }

  update() {
    fetch("http://localhost:8080/route/" + this.props.route)
      .then(response => response.json())
      .then(json => {
        if (json.length > 0) {
          let minTime = moment(json[0].timestamp).unix();
          let maxTime = moment(json[0].timestamp).unix();
          let minProg = json[0].progress;
          let maxProg = json[0].progress;
          console.log(minTime);

          for (let i = 0; i < json.length; i++) {
            if (moment(json[i].timestamp).unix() < minTime) {
              minTime = moment(json[i].timestamp).unix();
            }
            if (moment(json[i].timestamp).unix() > maxTime) {
              maxTime = moment(json[i].timestamp).unix();
            }
            if (json[i].progress < minProg) {
              minProg = json[i].progress;
            }
            if (json[i].progress > maxProg) {
              maxProg = json[i].progress;
            }
          }

          this.setState({
            data: json,
            maxTime: maxTime,
            minTime: minTime,
            minProg: minProg,
            maxProg: maxProg
          });
        }
      });
  }

  pointToString = (p, t) => {
    return this.progToScreen(p) + " " + this.timeToScreen(t) + " ";
  };

  progToScreen = p => {
    return (
      (p - this.state.minProg) /
        (this.state.maxProg - this.state.minProg) *
        (this.props.width - this.state.margins[0] - this.state.margins[2]) +
      this.state.margins[0]
    );
  };

  timeToScreen = t => {
    return (
      (t - this.state.minTime) /
        (this.state.maxTime - this.state.minTime) *
        (this.props.height - this.state.margins[1] - this.state.margins[3]) +
      this.state.margins[1]
    );
  };

  render() {
    let paths = {};
    let dompaths = [];
    let done = [-1];

    for (let i = 0; i < this.state.data.length; i++) {
      let dists = done.map(d => {
        return Math.abs(d - this.state.data[i].progress);
      });
      let minval = Math.min(...dists);
      if (minval > 1.0) {
        dompaths.push(
          <g>
            <path
            className="prog-tick"
              d={
                "M" +
                this.pointToString(
                  this.state.data[i].progress,
                  this.state.minTime
                ) +
                "L" +
                this.pointToString(
                  this.state.data[i].progress,
                  this.state.maxTime
                )
              }
            />
            <text
              className="prog-label"
              transform={
                "rotate(45," +
                (this.progToScreen(this.state.data[i].progress) - 5) +
                "," +
                (this.props.height - this.state.margins[3] + 10) +
                ")"
              }
              x={this.progToScreen(this.state.data[i].progress) - 5 + "px"}
              y={this.props.height - this.state.margins[3] + 10 + "px"}
            >
              {this.state.data[i].stop.name}
            </text>
          </g>
        );

        done.push(this.state.data[i].progress);
        console.log(done);
      }

      if (!(this.state.data[i].trip.id in paths)) {
        paths[this.state.data[i].trip.id] =
          "M" +
          this.pointToString(
            this.state.data[i].progress,
            moment(this.state.data[i].timestamp).unix()
          );
      }

      paths[this.state.data[i].trip.id] +=
        "L" +
        this.pointToString(
          this.state.data[i].progress,
          moment(this.state.data[i].timestamp).unix()
        );
    }

    console.log(paths);

    dompaths = dompaths.concat(
      Object.values(paths).map(p => {
        return <path d={p} />;
      })
    );

    let t = moment.unix(this.state.minTime);
    t.add(1, "hour");
    t.minute(0);

    while (t < moment.unix(this.state.maxTime)) {
      console.log(this.timeToScreen(t.unix()));
      dompaths.push(
        <g className="time">
          <text
            className="time-label"
            x="5px"
            y={this.timeToScreen(t.unix()) + 6 + "px"}
          >
            {t.format("ddd, hA")}
          </text>
          <path
            className="time-tick"
            d={
              "M" +
              this.pointToString(this.state.minProg, t.unix()) +
              "L" +
              this.pointToString(this.state.maxProg, t.unix())
            }
          />
        </g>
      );
      t.add(1, "hour");
    }

    console.log(moment.unix(this.state.maxTime));
    console.log(moment.unix(this.state.minTime));
    console.log(
      moment.unix(this.state.maxTime) - moment.unix(this.state.minTime)
    );

    if (this.state.data.length > 0) {
      return (
        <svg width={this.props.width} height={this.props.height}>
          <g fill="none" stroke={"#" + this.state.data[0].trip.route.color}>
            {dompaths}
          </g>
        </svg>
      );
    }  else {
      return <div />
    }
  }
}

export default MareyDiagram;
