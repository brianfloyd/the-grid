

class DOMSETUP {

constructor(){
  this.getData = async (api,param) => {
    if(!param)param="";
    else{ param = param.replace(/\//g, '-');}
    try {
                const response = await fetch(`${api}${param}`);
                const data = await response.json();
                console.log(data.errorCode)
                if(data.errorCode){
                    if(data.errorCode==='NOT_FOUND'  || 'GENERIC_ERROR'){
                        let paramApi = api.split('/').pop();
                        
                       console.log(paramApi)
                        if(api==='/workout/view/date/')return {"id":1,"date":"2024-08-24","sets":[{"id":2,"exerciseId":2,"name":"Seated Curls","group":"BICEP","weight":40,"reps":90,"count":0},{"id":5,"exerciseId":17,"name":"Overhead seated press","group":"TRICEP","weight":40,"reps":90,"count":0},{"id":3,"exerciseId":16,"name":"Dumbbell Skull Crusher","group":"TRICEP","weight":40,"reps":90,"count":1},{"id":4,"exerciseId":16,"name":"Dumbbell Skull Crusher","group":"TRICEP","weight":40,"reps":90,"count":1},{"id":6,"exerciseId":2,"name":"Seated Curls","group":"BICEP","weight":40,"reps":90,"count":3},{"id":1,"exerciseId":1,"name":"Concentrated Curls","group":"BICEP","weight":40,"reps":90,"count":4}]}
                        if(api.includes('/exercise/view/group/')) {
                            return [{"id":1,"name":"Concentrated Curls","group":"BICEP"},{"id":2,"name":"Seated Curls","group":"BICEP"},{"id":3,"name":"Straight Bar Curls","group":"BICEP"},{"id":4,"name":"Head Curls","group":"BICEP"},{"id":5,"name":"Seated Cable Curls","group":"BICEP"}];
                        }
                        if(api==='/exercise/view/groups/all')return [{"name":"BICEP","image":"images/icons/bicep.png","order":1},{"name":"BACK","image":"images/icons/back.png","order":2},{"name":"TRICEP","image":"images/icons/tricep.png","order":3},{"name":"CHEST","image":"images/icons/chest.png","order":4},{"name":"SHOULDER","image":"images/icons/shoulder.png","order":5},{"name":"LEGS","image":"images/icons/legs.png","order":6},{"name":"ABS","image":"images/icons/abs.png","order":7},{"name":"CARDIO","image":"images/icons/cardio.png","order":8},{"name":"MISC","image":"images/icons/misc.png","order":9}];
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
createGrid(expandGroup){
    let table =document.getElementById(`table-${expandGroup}`)
    if(expandGroup && table){
            while (table.firstChild) {
                table.removeChild(table.firstChild);
            }
        this.getData(`/exercise/view/group/${expandGroup}`,null).then(data=>{
         this.createGridItems(data);})
    }else{  
        const root =  document.getElementById('root');
        if(document.getElementsByClassName('grid-icons').length<1){
         const icondiv =  this.createElement('div',null,['grid-icons']);
            this.getData(`/exercise/view/groups/all`,null).then(groups=>{
                for (let muscleGroup of groups){
                icondiv.append(this.createElement('i',`icon-${muscleGroup.name.toLowerCase()}`,['icon']));
                root.append(icondiv);
                let img = this.createElement('img');
                img.src=muscleGroup.image.toLowerCase();
                let icon = document.getElementById(`icon-${muscleGroup.name.toLowerCase()}`)
                icon.appendChild(img);
                icon.addEventListener('click',function(){
                start.iconSelected(this);
                })
                }
            })
         this.createGridItems();
        }
    }
}
createGridItems(muscleGroup){
    start.getData(`/workout/view/date/`,start.getToday()).then(excerciseData=>{
        let legacyData={};
            if(muscleGroup){
             legacyData=excerciseData.sets;
             excerciseData.sets=muscleGroup;
            }
        for (let excercise of excerciseData.sets ){
            if(!document.getElementById(`table-${excercise.group}`)){
                const obj =start.createElement('div');
                const table=start.createElement('table');
                const tablebody=start.createElement('tbody');
                tablebody.id=`table-${excercise.group}`;
                document.getElementById(`icon-${excercise.group.toLowerCase()}`).classList.add('highlighted')
                table.append(tablebody);
                obj.append(table);
                obj.classList.add('the-grid');
                const grid = document.getElementById('the-grid');
                root.append(obj);
            }else{
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
                buttonElement.onclick = () => start.incrementExcercise(buttonElement, excercise, excerciseData.id);
                nextCell.appendChild(buttonElement);
                nextCell.id='button'
            if(muscleGroup){  
                    const rows = document.querySelectorAll('tr');
                    rows.forEach(row => {
                        // Get all <td> elements within the current row
                        const cells = row.querySelectorAll('td');
                        // Iterate backward to avoid index issues when removing elements
                        for (let i = cells.length - 1; i >= 0; i--) {
                            const cell = cells[i];
                            // Skip the first <td> and the <td> with id="button"
                            if (i !== 0 && cell.id !== 'button') {
                                row.deleteCell(i);
                            }
                        }
                    });
                    for(let ex of legacyData){
                
                        let buttons = document.querySelectorAll(`#table-${excerciseData.sets[0].group} .add-btn`);
                        buttons.forEach(function(button) {
                            const newButton = button.cloneNode(true);
                            if(ex.name ===excercise.name){
                            newButton.classList.add('included');
                            newButton.innerText='-';
                            }
                            button.parentNode.replaceChild(newButton, button);
                            newButton.addEventListener('click',function(){
                                const excerciseToAdd=this.parentElement.parentElement.firstChild.textContent;
                                const group = excerciseData.sets[0].group;
                                console.log(excerciseToAdd, group)


                            })
                        });
                        if(ex.name===excercise.name){
                            excerciseCell.classList.add('white-text');
                        }
                    }
                }
         }
                    
    }
             
    
        

    })
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
    console.log(excercise)

    fetch(`/set/view/${data.id}/save`, {
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
            elements.forEach(element=>this.clearTable(element))
            this.createGridItems();
        }else{
            let siblings = this.getSiblings(element);
            siblings.forEach(sibling=>sibling.classList.remove('edit-icon'))
            element.classList.add('edit-icon');
            let muscleGroup=element.id.split('-')[1].toUpperCase();
            this.createGrid(muscleGroup)
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
}
const start = new DOMSETUP();
start.createGrid()



