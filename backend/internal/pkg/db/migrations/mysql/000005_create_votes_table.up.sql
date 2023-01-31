CREATE TABLE IF NOT EXISTS Votes(
  ID INT NOT NULL UNIQUE AUTO_INCREMENT,
  UserID INT,
  LinkID INT,
  FOREIGN KEY (UserID) REFERENCES Users(ID),
  FOREIGN KEY (LinkID) REFERENCES Links(ID),
  PRIMARY KEY (ID)
)