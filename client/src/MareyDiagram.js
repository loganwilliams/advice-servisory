import React, { Component } from "react";
import moment from "moment";

class MareyDiagram extends Component {
  constructor(props) {
    super(props);
    this.state = {
      data: []
    };
  }

  componentDidMount() {
    fetch("http://localhost:8080/route/N")
      .then(response => response.json())
      .then(json => {
        let minTime = moment(json[0].timestamp).unix();
        let maxTime = moment(json[0].timestamp).unix();
        let minProg = json[0].progress;
        let maxProg = json[0].progress;

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

      });
  }

  render() {
    let paths = {};
    let dompaths = [];
    let width = 600;
    let height = 1000;

    for (let i = 0; i < this.state.data.length; i++) {
      if (!(this.state.data[i].trip.id in paths)) {
        paths[this.state.data[i].trip.id] =
          "M" +
          (this.state.data[i].progress - this.state.minProg) /
            (this.state.maxProg - this.state.minProg) *
            width +
          " " +
          (moment(this.state.data[i].timestamp).unix() - this.state.minTime) /
            (this.state.maxTime - this.state.minTime) *
            height +
          " ";
      }

      paths[this.state.data[i].trip.id] +=
        "L" +
        (this.state.data[i].progress - this.state.minProg) /
          (this.state.maxProg - this.state.minProg) *
          width +
        " " +
        (moment(this.state.data[i].timestamp).unix() - this.state.minTime) /
          (this.state.maxTime - this.state.minTime) *
          height +
        " ";
    }

    console.log(paths);

    dompaths = Object.values(paths).map(p => {
      return <path d={p} />;
    });

    return (
      <svg width={width} height={height}>
        <g fill="none" stroke="black">
          {dompaths}
        </g>
      </svg>
    );
  }
}

export default MareyDiagram;
