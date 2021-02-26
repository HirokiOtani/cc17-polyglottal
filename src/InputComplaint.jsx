//import React, { useEffect, useState } from "react";
//import axios from "axios";

export const InputComplaint = () => {
  return (
    <div>
      <form action="/addcomplaint" method="POST">
        {/* <input type="text" name="name" placeholder="name"></input> */}
        <input type="text" name="complaint" placeholder="complaint"></input>
        <button type="submit">Submit!</button>
      </form>
    </div>
  );
};
