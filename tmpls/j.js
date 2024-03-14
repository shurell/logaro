function addProg() {
    fetch("/registerprog", {
      method: "POST",
      body: JSON.stringify({id: 0, tname: sname.value, apitoken: stoken.value }),
      headers: {    "Content-type": "application/json; charset=UTF-8"  }
    });
    setTimeout(() => { // прячем через три секунды
        location.reload();
      }, 1500);
    }

function removeProg(el, id,at) {
    fetch("/removeprog", {
        method: "POST",
        body: JSON.stringify({id: id, tname: "", apitoken: at }),
      headers: {    "Content-type": "application/json; charset=UTF-8"  }
    });
    let elem = el.parentNode;
    elem.remove();
}

function addType() {
    fetch("/registertype", {
      method: "POST",
      body: JSON.stringify({id: 0, tname: tname.value, hexcolor: ttoken.value }),
      headers: {    "Content-type": "application/json; charset=UTF-8"  }
    });
    setTimeout(() => { // прячем через три секунды
        location.reload();
      }, 1500);
    }

function removeType(el, id, tn) {
    fetch("/removetype", {
        method: "POST",
        body: JSON.stringify({id: id, tname: tn, hexcolor: "" }),
      headers: {    "Content-type": "application/json; charset=UTF-8"  }
    });
    let elem = el.parentNode;
    elem.remove();
}