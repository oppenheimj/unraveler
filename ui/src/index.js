import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';

class Canvas extends React.Component {
  constructor(props) {
    super(props)
    this.state = { nodes: this.props.nodes, edges: this.props.edges }
  }

  componentWillReceiveProps(nextProps) {
    this.setState({ nodes: nextProps.nodes, edges: nextProps.edges })
  }

  edgesToArray() {
    var arr = []

    this.state.edges.forEach((toArray, from) => toArray.forEach(to => arr.push([from, to])))

    return arr
  }

  render() {
    const dim = 400
    return (
      <svg width={dim} height={dim}>
        {this.state.nodes.map(
          node => <circle cx={node[0]+dim/2} cy={node[1]+dim/2} r="3" stroke="black" stroke-width="4" fill="black" />
        )}
        {this.edgesToArray().map(([from, to]) =>
          <line
            x1={this.state.nodes[from][0]+dim/2}
            y1={this.state.nodes[from][1]+dim/2}
            x2={this.state.nodes[to][0]+dim/2}
            y2={this.state.nodes[to][1]+dim/2}
            stroke="black"
            strokeWidth={1}
          />
        )}
        {this.state.nodes.filter((node, i) => i % 3 === 0 && i !== this.state.nodes.length-2)
        .map(node => <text fontSize="8" x={node[0]+dim/2-3} y={node[1]+dim/2+3} fill="red">C</text>
        )}
        {this.state.nodes.filter((node, i) => i % 3 !== 0 || i === this.state.nodes.length-2)
        .map(node => <text fontSize="8" x={node[0]+dim/2-3} y={node[1]+dim/2+3} fill="yellow">H</text>
        )}
      </svg>
    )
  }
}

class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = { nodes: [], edges: [], i: 0, err: 0 }

    
    this.initWebsocket = this.initWebsocket.bind(this)
    this.initWebsocket()

    this.update = this.update.bind(this)
  }

  initWebsocket() {
    // Create WebSocket connection.
    this.socket = new WebSocket('ws://localhost:8080');

    // Listen for messages
    this.socket.addEventListener('message', event => {
      const data = JSON.parse(event.data);
      this.setState(data);
    });
  }

  update(e) {
    this.socket.send('Another one!');
  }

  render() {
    return (
      <div>
        <button onClick={this.update}>Draw graph</button>
        <p>i: {this.state.i}, err: {this.state.err}</p>
        <Canvas nodes={this.state.nodes} edges={this.state.edges} />
      </div>
    )
  }
}

ReactDOM.render(
  <App />,
  document.getElementById('root')
);
