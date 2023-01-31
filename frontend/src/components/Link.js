import React from "react";
import { AUTH_TOKEN } from "../constants";
import { timeDifferenceForDate } from "../utils";

const Link = (props) => {
  const { link } = props;
  const authToken = localStorage.getItem(AUTH_TOKEN);
  return (
    <div className="flex mt2 items-start">
      <div className="flex items-center">
        <span className="gray">{props.index + 1}.</span>
        {authToken && (
          <div
            className="ml1 gray f11"
            style={{ cursor: "pointer" }}
            onClick={() => {
              console.log("Clicked vote button");
            }}
          >
            ▲
          </div>
        )}
      </div>
      <div className="ml1">
        <div>
          {link.description} ({link.url})
        </div>
        {
          <div className="f6 1h-copy gray">
            {link.votes.length} votes | by{" "}
            {link.postedBy ? link.postedBy.username : "Anonymous"}{" "}
            {timeDifferenceForDate(link.createdAt)}
          </div>
        }
      </div>
    </div>
  );
};

export default Link;
