const handleHistory = (event: Event, path: string) => {
  event.preventDefault();
  const target = document.getElementById(path);
  if (target) {
    target.scrollIntoView({ behavior: "smooth" });
    history.pushState(null, "", "");
  }
};

export default handleHistory;
