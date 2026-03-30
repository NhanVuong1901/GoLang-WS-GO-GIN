let ws = null;
let localStream = null;
let peerConnection = null;
let currentTargetId = "";

const $ = (id) => document.getElementById(id);

const myIdInput = $("myId");
const targetIdInput = $("targetId");
const connectBtn = $("connectBtn");
const startCameraBtn = $("startCameraBtn");
const callBtn = $("callBtn");
const hangupBtn = $("hangupBtn");
const localVideo = $("localVideo");
const remoteVideo = $("remoteVideo");
const logBox = $("log");
