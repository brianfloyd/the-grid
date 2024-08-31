

class DOMSETUP {

constructor(){
  this.getData = async (api,param) => {
    if(!param)param="";
    else{ param = param.replace(/\//g, '-');}
    try {
                const response = await fetch(`${api}${param}`);
                const data = await response.json();
            
                if(data.errorCode){
                    if(data.errorCode==='NOT_FOUND'  || 'GENERIC_ERROR'){
                        let paramApi = api.split('/').pop();
                        return 'Date has no data'
                    //    console.log(paramApi)
                    //     if(api==='/workout/view/date/')return {"id":1,"date":"2024-08-24","sets":[{"id":2,"exerciseId":2,"name":"Seated Curls","group":"BICEP","weight":40,"reps":90,"count":0},{"id":5,"exerciseId":17,"name":"Overhead seated press","group":"TRICEP","weight":40,"reps":90,"count":0},{"id":3,"exerciseId":16,"name":"Dumbbell Skull Crusher","group":"TRICEP","weight":40,"reps":90,"count":1},{"id":4,"exerciseId":16,"name":"Dumbbell Skull Crusher","group":"TRICEP","weight":40,"reps":90,"count":1},{"id":6,"exerciseId":2,"name":"Seated Curls","group":"BICEP","weight":40,"reps":90,"count":3},{"id":1,"exerciseId":1,"name":"Concentrated Curls","group":"BICEP","weight":40,"reps":90,"count":4}]}
                    //     if(api.includes('/exercise/view/group/')) {
                    //         return [{"id":1,"name":"Concentrated Curls","group":"BICEP"},{"id":2,"name":"Seated Curls","group":"BICEP"},{"id":3,"name":"Straight Bar Curls","group":"BICEP"},{"id":4,"name":"Head Curls","group":"BICEP"},{"id":5,"name":"Seated Cable Curls","group":"BICEP"}];
                    //     }
                    //     if(api==='/exercise/view/groups/all')return [{"name":"BICEP","image":"images/icons/bicep.png","order":1},{"name":"BACK","image":"images/icons/back.png","order":2},{"name":"TRICEP","image":"images/icons/tricep.png","order":3},{"name":"CHEST","image":"images/icons/chest.png","order":4},{"name":"SHOULDER","image":"images/icons/shoulder.png","order":5},{"name":"LEGS","image":"images/icons/legs.png","order":6},{"name":"ABS","image":"images/icons/abs.png","order":7},{"name":"CARDIO","image":"images/icons/cardio.png","order":8},{"name":"MISC","image":"images/icons/misc.png","order":9}];
                    }
                }
               

                return data;
            } catch (error) {
                console.error('There has been a problem with your fetch operation:', error);
            }
    }

    this.registerServiceWorker =(()=>{
        if ('serviceWorker' in navigator) {
            navigator.serviceWorker.register('/service-worker.js')
              .then((registration) => {
                console.log('Service Worker registered with scope:', registration.scope);
              })
              .catch((error) => {
                console.log('Service Worker registration failed:', error);
              });
          }

    })();

}

async createGrid2(){


        const allExerciseGroups = await this.getData(`/exercise/view/groups/all`,null);
        if(!allExerciseGroups){
            console.log('no data')
            return;
        }
        const root =  document.getElementById('root');
        const icondiv =  this.createElement('div',null,['grid-icons']);
            for (const group of allExerciseGroups){
                root.append(icondiv);
               //createIcons
               const img = this.createElement('img');
                img.src=group.image.toLowerCase();
                icondiv.append(this.createElement('i',`icon-${group.name.toLowerCase()}`,['icon']));
                const icon = document.getElementById(`icon-${group.name.toLowerCase()}`);
                icon.appendChild(img);
                root.append(icondiv);
                icon.addEventListener('click',function(){
                start.iconSelected(this);
                 //create containers
                 
               
            })
        }
        for (const group of allExerciseGroups){
            const obj =start.createElement('div');
            const table=start.createElement('table');
            const tablebody=start.createElement('tbody');
            tablebody.id=`table-${group.name}`;
            table.append(tablebody);
            obj.append(table);
            obj.classList.add('the-grid');
            const grid = document.getElementById('the-grid');
            root.append(obj);
        }
        const noData = start.createElement('div','no-data',['hide']);
        noData.textContent= 'No data yet for this date';
        root.append(noData)
           
}

async createThisDate(){
    const thisDateData= await this.getData(`/workout/view/date/`,this.workoutDate);
    if(thisDateData==='Date has no data'){
        document.getElementById('no-data').classList.remove('hide');
        return;

    }else{
        document.getElementById('no-data').classList.add('hide')
    }
    for(const excercise of thisDateData.sets){
       
        document.getElementById(`icon-${excercise.group.toLowerCase()}`).classList.add('highlighted')
    const tb = document.getElementById(`table-${excercise.group}`);
    const row = tb.insertRow();
    const excerciseCell = row.insertCell(0);
    const weightCell = row.insertCell(1);
    const repsCell = row.insertCell(2);
    const nextCell = row.insertCell(3);
    const count = row.insertCell(4);
    count.setAttribute('onclick', 'start.decrementExcercise(this)');
    count.textContent=excercise.count;
if(excercise.count===0)count.textContent='';
    excerciseCell.textContent = excercise.name || '';
    weightCell.classList.add('data-cell','updatable');
    repsCell.classList.add('data-cell','updatable');
    weightCell.setAttribute('onclick', 'start.makeEditable(this)');
    repsCell.setAttribute('onclick', 'start.makeEditable(this)');
    weightCell.textContent = excercise.weight;
    repsCell.textContent = excercise.reps;
    const buttonElement = start.createElement('button', null, ['add-btn']);
    buttonElement.innerText = '+';
    buttonElement.onclick = () => start.incrementExcercise(buttonElement, excercise, thisDateData.id);
    nextCell.appendChild(buttonElement);
    nextCell.id='button'
    }
}

async createGridDayEditor(muscleGroup){
    document.getElementById('no-data').classList.add('hide');
   const  exerciseToDisplay = await this.getData(`/exercise/view/group/${muscleGroup}`);
   console.log({exerciseToDisplay})
   const exerciseMap={};
   exerciseToDisplay.forEach(item=>exerciseMap[item.name]={id:item.id,group:item.group,name:item.name})
   const  thisDayExercise= await this.getData(`/workout/view/date/`,this.workoutDate);
   console.log(exerciseMap)
   const allGridItems = document.querySelectorAll('[id^="table-"]');
   Array.from(allGridItems).forEach(item => {
    while (item.firstChild) {
        item.removeChild(item.firstChild);
    }
    });
   for(const exercise of exerciseToDisplay){
 
    const tb = document.getElementById(`table-${exercise.group}`);
    const row = tb.insertRow();
    const exerciseCell = row.insertCell(0);
    const nextCell = row.insertCell(1);
    const buttonElement = start.createElement('button', null, ['add-btn']);
    exerciseCell.textContent = exercise.name;
    buttonElement.innerText = '+';
    buttonElement.onclick = () => start.editDayExercise('',exerciseMap[exercise.name],true)
    nextCell.appendChild(buttonElement);
    nextCell.id='button'
   
   } 
   if(thisDayExercise==='Date has no data')return;
   for(let ex of thisDayExercise.sets){
    let buttons = document.querySelectorAll(`#table-${thisDayExercise.sets[0].group} .add-btn`);
    buttons.forEach(function(button){
        let exerciseCell=button.parentNode.parentNode.firstChild;
        let exercise = exerciseCell.textContent;
        const newButton = button.cloneNode(true);
        if(ex.name ===exercise){
        exerciseCell.classList.add('white-text');
        newButton.classList.add('included');
        newButton.innerText='-';
        }
        button.parentNode.replaceChild(newButton, button);
        newButton.addEventListener('click',function(){
            const exerciseToAdd=exercise;
            let fn=false;

           if(this.textContent==='+')fn=true;
        console.log(exerciseToAdd,exerciseMap[exerciseToAdd])
           start.editDayExercise(thisDayExercise.id  || "",exerciseMap[exerciseToAdd],fn)


       })
    });
}
}

editDayExercise(workoutId,exerciseDetails,add){
    console.log(workoutId, exerciseDetails,add);
    fetch(`/set/view/save`, {
        method: "POST",
        headers: {
            'Content-Type': "application/json"
        },
        body: JSON.stringify({
            "workoutId":workoutId,
           "exerciseId":exerciseDetails.id,
           "name":exerciseDetails.name,
           "group":exerciseDetails.group,
           "date":this.workoutDate
          
        })
      
    }).then(response => {
        if (!response.ok) {
            // If the response has a status code outside the range of 2xx, throw an error
            return response.json().then(errorData => {
                console.error('Error Response:', errorData); // Log the error details from the server
                throw new Error('Server responded with an error!');
            });
        }
        return response.json(); // If response is okay, parse it as JSON
    })
    .then(data => {
        console.log('Success:', data); // Log the successful response data
    })
    .catch(err => {
        console.error('Fetch Error:', err); // Log any error that occurred during the fetch or response handling
    });
    

}


incrementExcercise(element,data, workoutId){
 
    let excercise=element.parentElement.parentElement.childNodes[0].innerText;
    let weight=element.parentElement.parentElement.childNodes[1].innerText;
    let reps=element.parentElement.parentElement.childNodes[2].innerText;
    let incrementColumn = element.parentElement.nextSibling;
    incrementColumn.classList.remove('hide-count');
    let currentValue = incrementColumn.textContent;
    if (currentValue =>1)currentValue++;
    if (currentValue==='')currentValue=1;
    incrementColumn.textContent=currentValue;
    console.log(data)

    fetch(`/set/view/${Number(data.id)}/save`, {
        method: "POST",
        headers: {
            'Content-Type': "application/json"
        },
        body: JSON.stringify({
            ...data,
            workoutId,
            count: currentValue
        })
      
    });
    console.log(JSON.stringify({
        ...data,
        workoutId,
        count: currentValue
      
    }))
  
}
decrementExcercise(element){
        console.log(element)
        let excercise = element.parentElement.childNodes[0].innerText;
        let incrementColumn = element;
        let currentValue = incrementColumn.textContent;
            if (currentValue =>1)currentValue = currentValue -1
            if (currentValue==='' || currentValue ===0){
                currentValue='';
                element.classList.add('hide-count')
            }
        incrementColumn.textContent=currentValue;
    };
makeEditable(td) {
            const originalValue = td.innerText;
            td.innerHTML = `<input type='text' value='${originalValue}' onblur='start.saveValue(this)' onkeydown='start.handleKey(event, this)' />`;
            const input = td.querySelector('input');
            input.focus();
            input.select();
        }
iconSelected(element){

        if(element.classList.contains('edit-icon')){
            console.log('already selected');
            element.classList.remove('edit-icon');
            let elements = document.querySelectorAll('[id^="table-"]');
            elements.forEach(element=>this.clearTable(element));
            this.createThisDate();
     
        }else{
            let siblings = this.getSiblings(element);
            siblings.forEach(sibling=>sibling.classList.remove('edit-icon'))
            element.classList.add('edit-icon');
            let muscleGroup=element.id.split('-')[1].toUpperCase();
            this.createGridDayEditor(muscleGroup)
        }


    }

createElement(kind,id,classList){
    const doc = document.createElement(kind);
    if(id)doc.id=id;
    if(Array.isArray(classList)){
        classList.forEach(className=>
            doc.classList.add(className)
        )
    }
    return doc;
}
saveValue(input) {
    console.log(input)
    const newValue = input.value;
    const td = input;
    td.innerHTML = newValue;
    td.onclick = function() { start.makeEditable(td); };
}

handleKey(event, input) {
    if (event.key === 'Enter') {
        this.saveValue(input);
    } else if (event.key === 'Escape') {
        const td = input.parentElement;
        td.innerHTML = input.value;
        td.onclick = function() { start.makeEditable(td); };
    }
}
getToday(){
        const date = new Date();
        let year = date.getFullYear();
        let month = (date.getMonth() + 1).toString().padStart(2, '0'); // Months are zero-based, so add 1
        let day = date.getDate().toString().padStart(2, '0');
        return `${year}/${month}/${day}`;
}
getSiblings(element) {
    // Get the parent node of the element
    let parent = element.parentNode;

    // Get all children of the parent node
    let children = Array.prototype.slice.call(parent.children);

    // Filter out the original element to get only siblings
    return children.filter(function(child) {
        return child !== element;
    });
}
clearTable(table){
    while (table.firstChild) {
        table.removeChild(table.firstChild);
    }
}
updateDate(){
    if(!this.currentDate)this.currentDate = new Date();
    let month = (this.currentDate.getMonth() + 1).toString(); // Months are zero-based, so add 1
    let day = this.currentDate.getDate().toString();
    let year = this.currentDate.getFullYear();
    document.getElementById('current-date').textContent=`${month}/${day}/${year}`;
    this.workoutDate=`${month}/${day}/${year}`;
    let elements = document.querySelectorAll('[id^="table-"]');
    elements.forEach(element=>this.clearTable(element));
    this.createThisDate()
}
addDateListeners(){
  
document.getElementById('prev-date').addEventListener('click', ()=>{
   
    this.currentDate.setDate(this.currentDate.getDate() - 1);
    this.updateDate();

})
document.getElementById('next-date').addEventListener('click', ()=>{
   
    this.currentDate.setDate(this.currentDate.getDate() + 1);
    this.updateDate();

})
document.getElementById('current-date').addEventListener('click', ()=>{
   
    this.currentDate=null;
    this.updateDate();

})


}
}
const start = new DOMSETUP();
//start.createGrid();
start.updateDate();
start.addDateListeners();
start.createGrid2();
start.createThisDate();


