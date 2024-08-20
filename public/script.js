

class DOMSETUP {

constructor(){
  this.getData= async (api,param) => {
    param = param.replace(/\//g, '-');
    try {
                const response = await fetch(`${api}${param}`);
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data = await response.json();
                return data;
            } catch (error) {
                console.error('There has been a problem with your fetch operation:', error);
            }
    
  }

}
createGrid(){
    this.getData(``).then(group=>{
    const icondiv =  this.createElement('div','icon-set',['grid-icons']);
    icondiv.append(this.createElement('i',`icon-${excercise.group}`));
    let img = this.createElement('img');
    img.src=
    document.getElementById(`icon-${excercise.group}`).appendChild(img)
    })

    this.getData(`/workout/view/date/`,this.getToday()).then(excerciseData=>{
    const root =  document.getElementById('root');
   
    root.append(icondiv);
        for (let excercise of excerciseData.sets ){
            if(!document.getElementById(`table-${excercise.group}`)){
            
                const obj =this.createElement('div');
                const table=this.createElement('table');
                const tablebody=this.createElement('tbody');
                tablebody.id=`table-${excercise.group}`;
                table.append(tablebody);
                obj.append(table);
                obj.classList.add('the-grid');
                const grid = document.getElementById('the-grid');
                root.append(obj);
            }
                const tb = document.getElementById(`table-${excercise.group}`);
                const row = tb.insertRow();
                const excerciseCell = row.insertCell(0);
                const weightCell = row.insertCell(1);
                const repsCell = row.insertCell(2);
                const nextCell = row.insertCell(3);
                const count = row.insertCell(4);
                count.classList.add('hide-count');
                count.setAttribute('onclick', 'start.decrementExcercise(this)');
                excerciseCell.textContent = excercise.name || '';
                weightCell.classList.add('data-cell','updatable');
                repsCell.classList.add('data-cell','updatable');
                weightCell.setAttribute('onclick', 'start.makeEditable(this)');
                repsCell.setAttribute('onclick', 'start.makeEditable(this)');
                weightCell.textContent = excercise.weight;
                repsCell.textContent = excercise.reps;
                nextCell.innerHTML='<td class="button-cell"><button class="add-btn" onclick="start.incrementExcercise(this)">+</button></td>';
        }
    })
}
incrementExcercise(element){
          
              let excercise=element.parentElement.parentElement.childNodes[0].innerText;
              let weight=element.parentElement.parentElement.childNodes[1].innerText;
              let reps=element.parentElement.parentElement.childNodes[2].innerText;
              
              let incrementColumn = element.parentElement.nextSibling;
              incrementColumn.classList.remove('hide-count');
              let currentValue = incrementColumn.textContent;
              if (currentValue =>1)currentValue++;
              if (currentValue==='')currentValue=1;
              incrementColumn.textContent=currentValue;
           
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


}


const start = new DOMSETUP();

start.createGrid()






// let gridLog={};
// const greenCircle='🟢';
// const redCircle='🔴';
// function makeEditable(td) {
//     const originalValue = td.innerText;
//     td.innerHTML = `<input type='text' value='${originalValue}' onblur='saveValue(this)' onkeydown='handleKey(event, this)' />`;
//     const input = td.querySelector('input');
//     input.focus();
//     input.select();
// }

// function saveValue(input) {
//     const newValue = input.value;
//     const td = input.parentElement;
//     td.innerHTML = newValue;
//     td.onclick = function() { makeEditable(td); };
// }

// function handleKey(event, input) {
//     if (event.key === 'Enter') {
//         saveValue(input);
//     } else if (event.key === 'Escape') {
//         const td = input.parentElement;
//         td.innerHTML = input.value;
//         td.onclick = function() { makeEditable(td); };
//     }
// }
// function incrementExcercise(element){

//           let excercise=element.parentElement.parentElement.childNodes[0].innerText;
//           let weight=element.parentElement.parentElement.childNodes[1].innerText;
//           let reps=element.parentElement.parentElement.childNodes[2].innerText;
//           excerciseId = excercise.split(' ').join('-');
//           let incrementColumn = document.getElementById(excerciseId);
//           incrementColumn.classList.remove('hide-count');
//           let currentValue = incrementColumn.textContent;
//           if (currentValue =>1)currentValue++;
//           if (currentValue==='')currentValue=1;
//           incrementColumn.textContent=currentValue;
//           let group = element.parentElement.parentElement.parentElement.id;
//           group = group.split('-')[1];
        
//           let date = document.getElementsByClassName('date')[0].innerText;
//           if(!gridLog[date])gridLog={[date]:{[group]:[{[excercise]:[weight,reps]}]}}
//           else if (!gridLog[date][group])gridLog[date][group]=[{[excercise]:[weight,reps]}]
//           else gridLog[date][group].push({[excercise]:[weight,reps]})

//         save(gridLog)
// }
// function decrementExcercise(element){
//     let excercise = element.parentElement.childNodes[0].innerText;
//    let group = element.parentNode.parentNode;
//    group = group.split('-')[1];
//    console.log(group)

//         excercise = excercise.split(' ').join('-')
//     let incrementColumn = document.getElementById(excercise);
//     let currentValue = incrementColumn.textContent;
//         if (currentValue =>1)currentValue = currentValue -1
//         if (currentValue==='' || currentValue ===0){
//             currentValue='';
//             element.classList.add('hide-count')
//         }
//     incrementColumn.textContent=currentValue;
// };
// function save(gridLog){
//     fetch('/gridLog', {
//         method: 'POST',
//         headers: {
//             'Content-Type': 'application/json',
//         },
//         body: JSON.stringify({ gridLog }),
//     })
//     .then(response => response.json())
//     .then(data => {
//         console.log('Success:', data);
//     })
//     .catch((error) => {
//         console.error('Error:', error);
//     });
// }
// function arrange(){
//     document.getElementById('landing-listener').classList.add('hide');
//     let arrangeEl=document.getElementById('arrange-container');
//     arrangeEl.classList.remove('hide');
//     if(document.getElementById('arrange-container').children.length<2)
//     getData('groups').then((data)=>{
//         for(group of data){
//             let newDiv = document.createElement('div');
//             // Add class to the new div
//             newDiv.className = `arrange-${i} arrange-group`;
        
//             newDiv.id=`dropdown-btn-${group}`;
//             let span = document.createElement('span');
//             span.innerText=`${group} 🔻`;
//             span.id=`${group}-rollup`;
//             span.addEventListener('click',createDropdown)
//             newDiv.append(span)
//             // Append the new div to the container
//             let dateElement =document.getElementById('id-date');
//             let rootElement = document.getElementById('arrange-container');
//             rootElement.append(newDiv);
//             document.getElementsByTagName('body')[0].append(dateElement)
//         }
//     });

//     function createDropdown() {  
//         let parent= this.parentNode;
//         this.removeEventListener('click',createDropdown);
//         this.addEventListener('click',toggleDropdown);
//         let group = parent.id.split('-')[2];
//         let element = document.getElementById(parent.id);
//         fetchExerciseData().then(data=>{
//             data.forEach((ex)=>{
//                 if(ex.Category===group){
//                     let excercise = ex.excercise;
//                     let item = document.createElement('p');
//                     let span =document.createElement('span');
//                     span.className='span-check';
//                     span.innerText=redCircle;
//                     item.innerText = excercise;
//                     item.classList.add(`roll-down`,`${group}`);
//                     item.append(span);
//                     span.addEventListener('click',toggleCheck);
//                     element.append(item);
//                 }
//             })
//         });
//       }
//     function toggleDropdown(){
//         let a = this.parentNode.children;
//         for (let i = 0; i < a.length; i++) {
//             if(a[i].classList.contains('hide'))
//                 a[i].classList.remove('hide')
//             else if(a[i].tagName==='P')a[i].classList.add('hide');  
//     }}
//     function toggleCheck(){
//        if(this.innerText === greenCircle)this.innerText=redCircle;
//        else if(this.innerText ===redCircle)this.innerText=greenCircle;
//        let excercise = this.parentNode.innerText.slice(0,this.parentNode.innerText.length-2)
//        let date = document.getElementById('day-date').innerText;
//        console.log(excercise,date)
//     }
// }
// function attachMainScreenListeners(){
// document.getElementById('landing-listener').addEventListener('click', function(event) {
//     if (event.target && event.target.classList.contains('landing-buttons')) {
//         if(event.target.id==='Arrange') arrange();
//         if(event.target.id==='Work') {
//             document.getElementById(this.id).classList.add('hide');
//             document.getElementById('the-grid').classList.remove('hide');
//         }
//     }
// document.getElementById('title-button').addEventListener('click',()=>{
//     document.getElementById('the-grid').classList.add('hide');
//     document.getElementById('arrange-container').classList.add('hide');
//     document.getElementById('landing-listener').classList.remove('hide');
// })
// });
// let dayButtons= document.getElementsByClassName('day-button');
// let dayDate = document.getElementById('day-date');
// let parent;
// for (let i = 0; i < dayButtons.length; i++) {
 
//     if(dayButtons[i].innerText==='✓'){
            
//         return;
//     }else{
//     dayButtons[i].addEventListener('click', function() {
     
//        for(el of dayButtons){
//         el.classList.remove('green');
//        }
//        this.classList.add('green');
    
//       parent=this;
//       document.getElementById('arrange-date').classList.remove('hide');
    
//       dayDate.innerText=getNextDayDate(this.innerText);
//     });
// };
// }
//     let increment=0;
//     let left = document.getElementById('left');
//     let right = document.getElementById('right');
//     left.addEventListener('click',function(){
//         increment--;
//         if (increment <0)increment=0;
//         dayDate.innerText=getNextDayDate(parent.innerText,increment);
//     })
//     right.addEventListener('click',function(){
//         increment++;
//         console.log(this)
//         dayDate.innerText=getNextDayDate(parent.innerText,increment);
//     })

// }
// function makeEntry(day,date,excercise){

//     date=document.getElementById('');
//     day = document.getElementsByClassName('green');
//     console.log(day[0])


// }
// async function getData(kind) {
//     try {
//         const response = await fetch('/excercises');
//         if (!response.ok) {
//             throw new Error('Network response was not ok');
//         }
//         const data = await response.json();
//         if (kind==='groups'){
//             const muscleGroupList=[];
//             for (i of data){
//               if (!muscleGroupList.includes(i.Category)) muscleGroupList.push(i.Category);
//             }
//             return muscleGroupList;
//         }
//         else return data;
//     } catch (error) {
//         console.error('There has been a problem with your fetch operation:', error);
//     }
// }
// async function fetchExerciseData() {
//     try {
//         const response = await fetch('/excercises');
//         if (!response.ok) {
//             throw new Error('Network response was not ok');
//         }
//         const data = await response.json();
//         return data;
//     } catch (error) {
//         console.error('There has been a problem with your fetch operation:', error);
//     }
// }
// function getNextDayDate(dayOfWeek, weeksForward = 0) {
//     if(dayOfWeek==='✓')return
//     const dayMap = {
//         'M': 1,
//         'T': 2,
//         'W': 3,
//         'TH': 4,
//         'F': 5,
//         'SA': 6,
//         'SU': 0
//     };

//     const today = new Date();
//     const currentDay = today.getDay();
//     const targetDay = dayMap[dayOfWeek.toUpperCase()];

//     if (targetDay === undefined) {
//         throw new Error('Invalid day of the week provided.');
//     }

//     // Calculate days until the target day
//     let daysUntilTarget = targetDay - currentDay;
//     if (daysUntilTarget < 0) {
//         daysUntilTarget += 7;
//     }

//     // Add the additional weeks if the second parameter is provided
//     const totalDays = daysUntilTarget + (weeksForward * 7);

//     // Get the target date
//     const targetDate = new Date(today);
//     targetDate.setDate(today.getDate() + totalDays);

//     // Format the date as mm/dd/yy
//     const formattedDate = targetDate.toLocaleDateString('en-US', {
//         month: '2-digit',
//         day: '2-digit'
//     });

//     return formattedDate;
// }




// function init(data){
//         let date = new Date;
//         const dayOfWeek = date.getDay();
//         const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
//         const dayName = days[dayOfWeek];
//         let dateElement = document.createElement('div');
//         dateElement.classList.add('date');
//         dateElement.id=('id-date')
//         dateElement.textContent=`${dayName}    ${date.getMonth()+1}/${date.getDate()}/${date.getFullYear()}`
//         document.getElementById('root').append(dateElement)
//        attachMainScreenListeners();
//         if ('serviceWorker' in navigator) {
//             navigator.serviceWorker.register('/service-worker.js')
//               .then((registration) => {
//                 console.log('Service Worker registered with scope:', registration.scope);
//               })
//               .catch((error) => {
//                 console.log('Service Worker registration failed:', error);
//               });
//           }
//           populateTable(data);
    
//         function populateTable(data) {
//             const muscleGroupList=[];
//           for (i of data){
//             if (!muscleGroupList.includes(i.Category)) muscleGroupList.push(i.Category);
//           }
         
//           muscleGroupList.forEach((muscleGroup)=>{
//             let obj = document.createElement('div');
//             let icon =document.createElement('img');
//             let table=document.createElement('table');
//             let tablebody=document.createElement('tbody');
//             tablebody.id=`table-${muscleGroup}`;
//             table.append(tablebody);
//             icon.setAttribute('src',`images/${muscleGroup}.png`);
//             icon.classList.add('icon-placeholder');
//             obj.append(icon);
//             obj.append(table)
//             obj.id=muscleGroup;
//             obj.classList.add('the-grid')
//             let grid = document.getElementById('the-grid');

           
//             grid.append(obj);
//           })
//          data.forEach(item => {
//                 const tableBody = document.getElementById(`table-${item.Category}`)
//                 const row = tableBody.insertRow();
//                 const exerciseCell = row.insertCell(0);
//                 const weightCell = row.insertCell(1);
//                 const repsCell = row.insertCell(2);
//                 const nextCell = row.insertCell(3);
//                 const count = row.insertCell(4);
//                 count.id = item.excercise.split(' ').join('-');
//                 count.classList.add('hide-count');
//                 count.setAttribute('onclick', 'decrementExcercise(this)');
//                 exerciseCell.textContent = item.excercise || '';
//                 weightCell.classList.add('data-cell','updatable');
//                 repsCell.classList.add('data-cell','updatable');
//                 weightCell.setAttribute('onclick', 'makeEditable(this)');
//                 repsCell.setAttribute('onclick', 'makeEditable(this)');
//                 weightCell.textContent = item.weight || '';
//                 repsCell.textContent = item.rep || '';
//                 nextCell.innerHTML='<td class="button-cell"><button class="add-btn" onclick="incrementExcercise(this)">+</button></td>';
//             });
//         }
      

// }
// fetchExerciseData().then(data=>init(data))
