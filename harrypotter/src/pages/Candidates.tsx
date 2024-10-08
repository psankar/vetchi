import { Divider, Flex, Input, Table } from "antd";
import { useState } from "react";
import { Link } from "react-router-dom";
import t from "../i18n/i18n";
import { formInputStyle } from "../Styles";

const data = [
  {
    id: "cand0123123",
    name: "Jeff Dean",
    lastPosition: "Senior Fellow",
    lastCompany: "Google",
    shortlistedOpenings: [
      {
        id: "JAN14-1",
        hiringManager: "a@example.com",
        title: "Distinguished Engineer",
      },
      {
        id: "JAN14-2",
        hiringManager: "b@example.com",
        title: "Fellow",
      },
      {
        id: "JAN14-3",
        hiringManager: "c@example.com",
        title: "Senior Fellow",
      },
    ],
  },
  {
    id: "cand1263561",
    name: "Sanjay Ghemawat",
    lastPosition: "Senior Fellow",
    lastCompany: "Google",
    shortlistedOpenings: [
      {
        id: "JAN14-1",
        hiringManager: "a@example.com",
        title: "Distinguished Engineer",
      },
      {
        id: "JAN14-2",
        hiringManager: "b@example.com",
        title: "Fellow",
      },
      {
        id: "JAN14-3",
        hiringManager: "c@example.com",
        title: "Senior Fellow",
      },
    ],
  },
  {
    id: "cand74358738",
    name: "Swami Sivasubramanian",
    lastPosition: "Vice President and General Manager",
    lastCompany: "Amazon Web Services",
    shortlistedOpenings: [
      {
        id: "JAN14-1",
        hiringManager: "a@example.com",
        title: "Distinguished Engineer",
      },
      {
        id: "JAN14-2",
        hiringManager: "b@example.com",
        title: "Fellow",
      },
      {
        id: "JAN14-3",
        hiringManager: "c@example.com",
        title: "Senior Fellow",
      },
    ],
  },
  {
    id: "cand045646456",
    name: "Joydeep Sen Sarma",
    lastPosition: "CTO",
    lastCompany: "Clearfeed",
    shortlistedOpenings: [
      {
        id: "JAN14-1",
        hiringManager: "a@example.com",
        title: "Distinguished Engineer",
      },
      {
        id: "JAN14-2",
        hiringManager: "b@example.com",
        title: "Fellow",
      },
      {
        id: "JAN14-3",
        hiringManager: "c@example.com",
        title: "Senior Fellow",
      },
    ],
  },
  {
    id: "cand078979789",
    name: "Richard Hipp",
    lastPosition: "Maintainer",
    lastCompany: "Hwaci",
    shortlistedOpenings: [
      {
        id: "JAN14-1",
        hiringManager: "a@example.com",
        title: "Distinguished Engineer",
      },
      {
        id: "JAN14-2",
        hiringManager: "b@example.com",
        title: "Fellow",
      },
      {
        id: "JAN14-3",
        hiringManager: "c@example.com",
        title: "Senior Fellow",
      },
    ],
  },
];

export default function Candidates() {
  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
    },
    {
      title: "Last Position",
      dataIndex: "lastPosition",
      key: "lastPosition",
    },
    {
      title: "Last Company",
      dataIndex: "lastCompany",
      key: "lastCompany",
    },
    {
      title: "Shortlisted Openings",
      dataIndex: "shortlistedOpenings",
      key: "shortlistedOpenings",
    },
    {
      title: "Candidacy",
      key: "candidacy",
      render: (record: { id: string }) => (
        <Link to={`/candidacy/${record.id}`}>Candidacy</Link>
      ),
    },
  ];

  const dataSource = data.map((candidate, index) => ({
    key: index,
    ...candidate,
    shortlistedOpenings: (
      <Table
        columns={[
          { title: "Title", dataIndex: "title", key: "title" },
          {
            title: "Hiring Manager",
            dataIndex: "hiringManager",
            key: "hiringManager",
          },
        ]}
        dataSource={candidate.shortlistedOpenings}
        pagination={false}
        size="small"
        scroll={{ x: true }}
      />
    ),
  }));

  const [filterText, setFilterText] = useState("");

  return (
    <Flex vertical style={{ margin: "20px" }}>
      <Input
        type="text"
        placeholder={t("applications.filter_by_name_or_employer")}
        value={filterText}
        onChange={(e) => setFilterText(e.target.value)}
        style={formInputStyle}
      />
      <Divider />
      <Table columns={columns} dataSource={dataSource} scroll={{ x: true }} />;
    </Flex>
  );
}
