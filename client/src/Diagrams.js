import React, { Component } from "react";
import MareyDiagram from "./MareyDiagram.js";
import './Diagrams.css';

class Diagrams extends Component {
  render() {
    let lines = ["1", "2", "3", "4", "5", "6", "7", "A", "B", "C", "D", "F", "E","G", "J", "L", "M", "N", "Q", "R", "W"];
    let graphs = []

    for (let i = 0; i < lines.length; i++) {
      graphs.push(<div className="line-block">
        <div className="label">{lines[i] + " Train"}</div> 
        <MareyDiagram key={"MareyDiagram-" + lines[i]} width={500} height={800} route={lines[i]} direction={-1} /></div>);
    }

    return (
      <div className="graphs">
        {graphs}    
      </div>
      )
  }
}

export default Diagrams;
