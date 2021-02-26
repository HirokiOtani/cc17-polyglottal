import axios from "axios";

export default async function fetchComplaints() {
  const comps = await axios.get("/getcomplaints");
  return comps;
}
