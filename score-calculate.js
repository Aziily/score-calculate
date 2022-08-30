// import React from 'react';
// import ReactDOM from 'react-dom/client';
// import './score-Calculate.css';

$(document).ajaxError(function(){
    alert("请求失败，请确认查询信息！");
});

let allCredits = 0
let allScores = 0
let allGPAs = 0

function NewtonLoading() {
    return(
        <div className="cradle">
            <div className="dot"></div>
            <div className="dot"></div>
            <div className="dot"></div>
            <div className="dot"></div>
            <div className="dot"></div>
            <div className="dot"></div>
            <div className="dot"></div>
            <div className="dot"></div>
        </div>
    )
}

class One extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            choose: true,
            classname : this.props["data"][0],
            credit : this.props["data"][1],
            score : this.props["data"][2],
            gpa : this.props["data"][3],
        }
        this.handleChange = this.handleChange.bind(this)
    }
    componentDidMount() {
        let that = this
        allCredits += that.state.credit
        allScores += that.state.score * that.state.credit
        allGPAs += that.state.gpa * that.state.credit
    }
    handleChange() {
        // console.log("get")
        let that = this
        if (that.state.choose){
            allCredits -= that.state.credit
            allScores -= that.state.score * that.state.credit
            allGPAs -= that.state.gpa * that.state.credit
        } else {
            allCredits += that.state.credit
            allScores += that.state.score * that.state.credit
            allGPAs += that.state.gpa * that.state.credit
        }
        that.setState({
            choose: !that.state.choose
        })
        // console.log(allCredits, this.state.choose)
        ReactDOM.render(
            <MainBoard />,
            document.getElementById("main")
        )
    }
    render() {
        let that = this
        return (
            <div className="oneClass">
                <span className="choosebox"><button className="choosebtn" style={{"backgroundColor": that.state.choose ? "aqua" : "white"}} onClick={this.handleChange}></button></span>
                <span className="classname">{that.state.classname}</span>
                <span className="credit">{that.state.credit}</span>
                <span className="score">{that.state.score}</span>
                <span className="gpa">{that.state.gpa}</span>
            </div>
        )
    }

}

class Semester extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            Ones : [],
            isShow: false,
        }
        this.showMore = this.showMore.bind(this)
    }

    componentDidMount() {
        let that = this
        let newOnes = []
        for (let i = 0; i < that.props["data"]["scores"].length; i++) {
            if (that.props["data"]["scores"][i] == 0) continue
            newOnes.push(<One key={i} data={[that.props["data"]["classnames"][i], that.props["data"]["credits"][i], that.props["data"]["scores"][i], that.props["data"]["gpas"][i]]} />)
        }
        that.setState({
            Ones: newOnes
        })
    }
    showMore() {
        let that = this
        that.setState({
            isShow: !that.state.isShow
        })
    }

    render() {
        let that = this
        return (
            <div className='semester'>
                <div className="semTitle">
                    <img src="./src/triangle.png" className={that.state.isShow ? "showMore-Y" : "showMore"} onClick={that.showMore}></img>
                    <span>{that.props["year"]}</span>
                </div>
                <div style={{"display": that.state.isShow ? "flex" : "", "flexDirection": "column"}} className={that.state.isShow ? "appearMore" : "appearLess"}>
                    <div className="titles">
                        <span className="choosebox"></span>
                        <span className="classname">课程名称</span>
                        <span className="credit">学分</span>
                        <span className="score">成绩</span>
                        <span className="gpa">绩点</span>
                    </div>
                    {that.state.Ones}
                </div>
            </div>
        )
    }
}

class ShowBoard extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            Sems : []
        }
    }

    componentDidMount() {
        let that = this
        let newSems = []
        for(let year in that.props["data"]){
            newSems.push(<Semester key={year} data={that.props["data"][year]} year={year} />)
        }
        that.setState({
            Sems: newSems,
        })
    }

    render() {
        let that = this
        return (
            <div className='showBoard'>
                {that.state.Sems}
            </div>
        )
    }
}


class MainBoard extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            id: "",
            passwd: "",
            jsonMess : "",
            isLoading: false,
        }
        this.getMess = this.getMess.bind(this)
    }

    getMess() {
        // console.log(this.state)
        let that = this
        if (that.state.id.length != 10) {
            alert("请输入正确的学号！")
            return
        }
        that.setState({
            isLoading: true,
        })
        $.ajax({
            type: 'POST',
            url: "http://azily.natapp1.cc/api/search",
            data: JSON.stringify({
                id: that.state.id,
                passwd: window.btoa(that.state.passwd),
            }),
            dataType: 'json',
            success: function(data) {
                allCredits = 0
                allGPAs = 0
                allScores = 0
                that.setState({
                    jsonMess: data
                })
                ReactDOM.render(
                    <div></div>,
                    document.getElementById("show")
                )
                ReactDOM.render(
                    <ShowBoard data={data} />,
                    document.getElementById("show")
                )
                // // console.log(that.state.jsonMess)
                // for(let year in data){

                //     for (let i = 0; i < data[year]["scores"].length; i++) {
                //         if (parseInt(data[year]["scores"][i]) == 0) continue
                //         let credit = parseFloat(data[year]["credits"][i])
                //         let score = parseInt(data[year]["scores"][i])
                //         let gpa = parseFloat(data[year]["gpas"][i])
                //         allScores += score * credit
                //         allCredits += credit
                //         allGPAs += gpa * credit
                //     }
                // }
                that.setState({
                    isLoading: false,
                })
            },
            error: function (jqXHR, textStatus, errorThrown) {
                // console.log(jqXHR.responseText)
                that.setState({
                    isLoading: false,
                })
            }
        })
    }

    setID(e) {
        let value = e.target.value;
        this.setState({
            id: value
        })
    }
    setPasswd(e) {
        let value = e.target.value;
        this.setState({
            passwd: value
        })
        if (e.keyCode === 13) {
            this.getMess()
        }
    }

    render(){
        let that = this
        if (!that.state.isLoading){
            return(
                <div className="mainBoard">
                    <div className="login">
                        <input type="text" name="id" placeholder="studentid" value={that.state.id} className="inputBox" onChange={e => that.setID(e)}></input>
                        <input type="password" name="passwd" placeholder="passwd" value={that.state.passwd} className="inputBox" onChange={e => that.setPasswd(e)}></input>
                        <button className="btn-submit" onClick={that.getMess}>提交</button>
                    </div>
                    <div className="showAll">
                        <div id="averageGPA">平均绩点: {allCredits == 0 ? 5.00 : (allGPAs / allCredits).toFixed(2)}</div>
                        <div id="averageGPA">平均成绩: {allCredits == 0 ? 100.00 : (allScores / allCredits).toFixed(2)}</div>
                        <div id="allCredits">总学分: {allCredits}</div>
                    </div>
                </div>
            );
        } else {
            return(
                <div className="mainBoard">
                    <div className="login">
                        <input type="text" name="id" placeholder="studentid" value={that.state.id} className="inputBox" onChange={e => that.setID(e)}></input>
                        <input type="password" name="passwd" placeholder="passwd" value={that.state.passwd} className="inputBox" onChange={e => that.setPasswd(e)}></input>
                        <button className="btn-submit" onClick={that.getMess}>提交</button>
                    </div>
                    <NewtonLoading />
                </div>
            )
        }
    }

}

ReactDOM.render(
    <MainBoard />,
    document.getElementById("main")
)