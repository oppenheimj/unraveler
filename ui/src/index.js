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
    const dim = 2000
    return (
      <svg width={dim} height={dim}>
        {this.state.nodes.map(
          node => <circle cx={node[0]+dim/2} cy={node[1]+dim/2} r="3" stroke="black" strokeWidth="4" fill="black" />
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

class ParamForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      n: 100,
      kr: 100,
      ka: 1.0,
      kn: 1.0,
      maxIters: 1000,
      minError: 0.001,
      theta: 0.8,
      numThreads: 8
    };

    this.update = this.props.update.bind(this)

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    const { name, value } = event.target
    console.log("MATCH", value.match(/\.0*$/))

    if (!value || value.match(/\.0*$/)) {
      this.setState({ [name]: value });
    } else {
      const isInt = n => n % 1 === 0
      this.setState({ [name]: isInt(name) ? parseInt(value) : parseFloat(value) });
    }
  }

  handleSubmit(event) {
    event.preventDefault();
    this.update(JSON.stringify(this.state))
  }

  render() {
    return (
      <form onSubmit={this.handleSubmit}>
        <label>
          n:
          <input type="text" name='n' value={this.state.n} onChange={this.handleChange} />
        </label>
        <label>
          ka:
          <input type="text" name='ka' value={this.state.ka} onChange={this.handleChange} />
        </label>
        <label>
          kr:
          <input type="text" name='kr' value={this.state.kr} onChange={this.handleChange} />
        </label>
        <label>
          kn:
          <input type="text" name='kn' value={this.state.kn} onChange={this.handleChange} />
        </label>
        <label>
          maxIters:
          <input type="text" name='maxIters' value={this.state.maxIters} onChange={this.handleChange} />
        </label>
        <label>
          minError:
          <input type="text" name='minError' value={this.state.minError} onChange={this.handleChange} />
        </label>
        <label>
          theta:
          <input type="text" name='theta' value={this.state.theta} onChange={this.handleChange} />
        </label>
        <label>
          numThreads:
          <input type="text" name='numThreads' value={this.state.numThreads} onChange={this.handleChange} />
        </label>
        <input type="submit" value="Submit" />
      </form>
    );
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
    console.log('!!!', e)
    this.socket.send(e)
  }

  render() {
    return (
      <div>
        <ParamForm update={this.update}/>
        {/* <button onClick={this.update}>Draw graph</button> */}
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
