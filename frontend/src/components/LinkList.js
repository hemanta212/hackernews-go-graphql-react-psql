import { gql, useQuery } from "@apollo/client";
import React from "react";
import Link from "./Link";

const FEED_QUERY = gql`
  {
    feed {
      links {
        id
        createdAt
        description
        url
        postedBy {
          username
        }
        votes {
          id
        }
      }
    }
  }
`;

const LinkList = () => {
  const { data } = useQuery(FEED_QUERY);
  return (
    <div>
      {data && (
        <>
          {data.feed.links.map((link, index) => (
            <Link key={link.id} link={link} index={index} />
          ))}
        </>
      )}
    </div>
  );
};

export default LinkList;
