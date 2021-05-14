import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';

class Canvas extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      nodes: this.props.nodes,
      edges: this.props.edges,
      minX: this.props.minX,
      maxX: this.props.maxX,
      minY: this.props.minY,
      maxY: this.props.maxY
    }
  }

  componentWillReceiveProps(nextProps) {
    this.setState({
      nodes: nextProps.nodes,
      edges: nextProps.edges,
      minX: nextProps.minX,
      maxX: nextProps.maxX,
      minY: nextProps.minY,
      maxY: nextProps.maxY})
  }

  edgesToArray() {
    var arr = []

    this.state.edges.forEach((toArray, from) => toArray.forEach(to => arr.push([from, to])))

    return arr
  }

  render() {
    const originalWidth = this.state.maxX - this.state.minX
    const originalHeight = this.state.maxY - this.state.minY

    const dim = Math.max(originalWidth, originalHeight)

    const pointSize = 1 / Math.log10(this.state.nodes.length)

    const dim2Pct = (d, m) => `${(d - m)/dim * 100}%`

    return (
        <svg viewBox="0 0 100 100">
        {this.state.nodes.map(
          node => <circle cx={dim2Pct(node[0], this.state.minX)} cy={dim2Pct(node[1], this.state.minY)} r={pointSize} stroke="black" strokeWidth={0} fill="black" />
        )}
        {this.edgesToArray().map(([from, to]) =>
          <line
            x1={dim2Pct(this.state.nodes[from][0], this.state.minX)}
            y1={dim2Pct(this.state.nodes[from][1], this.state.minY)}
            x2={dim2Pct(this.state.nodes[to][0], this.state.minX)}
            y2={dim2Pct(this.state.nodes[to][1], this.state.minY)}
            stroke="black"
            strokeWidth={pointSize/6}
          />
        )}
        {/* {this.state.nodes.filter((node, i) => i % 3 === 0 && i !== this.state.nodes.length-2)
        .map(node => <text fontSize="8" x={node[0]+dim/2-3} y={node[1]+dim/2+3} fill="red">C</text>
        )}
        {this.state.nodes.filter((node, i) => i % 3 !== 0 || i === this.state.nodes.length-2)
        .map(node => <text fontSize="8" x={node[0]+dim/2-3} y={node[1]+dim/2+3} fill="yellow">H</text>
        )} */}
      </svg>
    )
  }
}

class ParamForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      n: 32,
      kr: 1000,
      ka: 1,
      kn: 1,
      maxIters: 1000,
      updateEvery: 10,
      minError: 0.0001,
      theta: 0,
      numThreads: 1
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
        </label><br />
        <label>
          ka:
          <input type="text" name='ka' value={this.state.ka} onChange={this.handleChange} />
        </label><br />
        <label>
          kr:
          <input type="text" name='kr' value={this.state.kr} onChange={this.handleChange} />
        </label><br />
        <label>
          kn:
          <input type="text" name='kn' value={this.state.kn} onChange={this.handleChange} />
        </label><br />
        <label>
          maxIters:
          <input type="text" name='maxIters' value={this.state.maxIters} onChange={this.handleChange} />
        </label><br />
        <label>
          updateEvery:
          <input type="text" name='updateEvery' value={this.state.updateEvery} onChange={this.handleChange} />
        </label><br />
        <label>
          minError:
          <input type="text" name='minError' value={this.state.minError} onChange={this.handleChange} />
        </label><br />
        <label>
          theta:
          <input type="text" name='theta' value={this.state.theta} onChange={this.handleChange} />
        </label><br />
        <label>
          numThreads:
          <input type="text" name='numThreads' value={this.state.numThreads} onChange={this.handleChange} />
        </label><br />
        <input type="submit" value="Submit" />
      </form>
    );
  }
}

class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = { nodes: [], edges: [], i: 0, err: 0, minX: 0, maxX: 0, minY: 0, maxY: 0 }

    
    this.initWebsocket = this.initWebsocket.bind(this)
    this.initWebsocket()

    this.update = this.update.bind(this)
  }

  initWebsocket() {
    // Create WebSocket connection.
    this.socket = new WebSocket('ws://localhost:8080');

    // Listen for messages
    this.socket.addEventListener('message', event => {
      this.setState(JSON.parse(event.data))
    });
  }

  update(e) {
    this.socket.send(e)
  }

  render() {
    const originalWidth = this.state.maxX - this.state.minX
    const originalHeight = this.state.maxY - this.state.minY
    const dim = Math.max(originalWidth, originalHeight)

    return (
      // width: {dim},
      <div>
        <ParamForm update={this.update}/>
        <p>i: {this.state.i}, err: {this.state.err}</p>
        <Canvas nodes={this.state.nodes} edges={this.state.edges} minX={this.state.minX} maxX={this.state.maxX} minY={this.state.minY} maxY={this.state.maxY} />
      </div>
    )
  }
}

ReactDOM.render(
  <App />,
  document.getElementById('root')
);
