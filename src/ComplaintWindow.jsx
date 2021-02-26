export const ComplaintWindow = (props) => {
  const complaint = props.complaint;
  return (
    <div style={{ position: "relative", bottom: 250, left: 650 }}>
      <table style={{ border: 3 }}>
        <tr>
          <th></th>
          <th>Score</th>
          <th>Magnitude</th>
        </tr>
        <tr>
          <th>Overall</th>
          <th>{complaint.Score}</th>
          <th>{complaint.Magnitude}</th>
        </tr>
        <tr>
          <th>{complaint.Word1}</th>
          <th>{complaint.Score1}</th>
          <th>{complaint.Magnitude1}</th>
        </tr>
        <tr>
          <th>{complaint.Word2}</th>
          <th>{complaint.Score2}</th>
          <th>{complaint.Magnitude2}</th>
        </tr>
        <tr>
          <th>{complaint.Word3}</th>
          <th>{complaint.Score3}</th>
          <th>{complaint.Magnitude3}</th>
        </tr>
      </table>
      <p
        style={{
          width: 200,
          height: 200,
        }}
      >
        {complaint.Text}
      </p>
    </div>
  );
};
