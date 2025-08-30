//  }
// }
// 
// // 保留原有的greet函数以保持向后兼容
// const nameElement = document.getElementById("name");
// if (nameElement instanceof HTMLElement) {
//   nameElement.focus();
// }
// const resultElement = document.getElementById("result");
// 
// // Setup the greet function
// window.greet = function () {
//     // Get name
//     let name = '';
//     if (nameElement instanceof HTMLInputElement) {
//       name = nameElement.value;
//     }
// 
//     // Check if the input is empty
//     if (name === "") return;
// 
//     // Call App.Greet(name)
//     try {
//         Greet(name)
//             .then((result: string) => {
//                 // Update result with data back from App.Greet()
//                 if (resultElement) {
//                   resultElement.innerText = result;
//                 }
//             })
//             .catch((err: Error) => {
//                 console.error(err);
//             });
//     } catch (err) {
//         console.error(err);
//     }
// };
