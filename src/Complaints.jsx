import React, { useEffect, useState } from "react";
import { InputComplaint } from "./InputComplaint";
import { ComplaintWindow } from "./ComplaintWindow";
import axios from "axios";
import "./Complaints.css";

export const Complaints = () => {
  const [complaints, setComplaints] = useState([
    [
      {
        ID: 1,
        CreatedAt: "2021-02-25T21:57:20.953305+09:00",
        UpdatedAt: "2021-02-25T21:57:20.953305+09:00",
        DeletedAt: null,
        Text:
          "2階から変な音が聞こえます。絶対に自分を監視している人がいるんです。これは幻覚なんかではありません。",
        PostedTime: "2021-02-24T00:00:00+09:00",
        PatientID: 1,
        Score: -0.3,
        Magnitude: 1,
        Word1: "音",
        Salience1: 0.56259376,
        Score1: -0.7,
        Magnitude1: 0.7,
        Word2: "幻覚",
        Salience2: 0.22695698,
        Score2: -0.2,
        Magnitude2: 0.2,
        Word3: "人",
        Salience3: 0.21044928,
        Score3: 0,
        Magnitude3: 0,
      },
    ],
    [],
    [],
    [],
    [],
    [],
  ]);
  const [patients, setPatients] = useState([{ Name: "Tarou", ID: 1 }]);
  const [complaint, setComplaint] = useState("complaint test");
  const [id, setId] = useState(1);
  useEffect(() => {
    async function fetchData() {
      try {
        const p = await axios.get("/getpatients");
        setPatients(p.data);
        const numPatient = patients.length;
        const compsPerPatient = [];
        const c = await axios.get("/getcomplaints");
        const comps = c.data;
        for (let i = 0; i < numPatient; i++) {
          const list = [];
          for (const comp of comps) {
            if (comp.PatientID - 1 === i) {
              list.push(comp);
            }
          }
          compsPerPatient.push(list);
        }

        setComplaints(compsPerPatient);
      } catch (err) {
        console.error(err);
      }
    }
    fetchData();
  }, [id]);
  //console.log(complaints);
  //console.log(patients);
  //console.log("AAAaaAAAA", complaints[id - 1]);
  let complaintList = <p>nothing</p>;

  if (complaints[id - 1]) {
    complaintList = complaints[id - 1].map((complaint) => {
      let cn = "normalGrid";
      let text = "";
      if (complaint.Magnitude > 1) cn = "alertGrid";
      if (
        complaint.Word1 === "幻覚" ||
        complaint.Word2 === "幻覚" ||
        complaint.Word3 === "幻覚"
      )
        text = "!";
      return (
        <button
          onClick={() => {
            setComplaint(complaint);
            alert(complaint.Text);
          }}
          className={cn}
        >
          {text}
        </button>
      );
    });
  }

  const patientList = patients.map((patient) => (
    <option value={patient.ID}>{patient.Name}</option>
  ));
  return (
    <div style={{ position: "relative", left: 10 }}>
      <form>
        <select id="patientSelect">{patientList}</select>
        <input
          type="button"
          value="select patient"
          onClick={() => {
            const n = document.getElementById("patientSelect").value;
            setId(parseInt(n));
            setComplaints(complaints);
          }}
        ></input>
      </form>
      <div className="container">{complaintList}</div>
      <ComplaintWindow complaint={complaint} />
    </div>
  );
};
