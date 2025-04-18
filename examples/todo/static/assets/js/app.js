var authToken = localStorage.getItem('auth_token');
var currentUser = null;

if (authToken) {
  // Try to parse stored user info
  try {
    currentUser = JSON.parse(localStorage.getItem('user_info'));
  } catch (e) {
    console.error('Failed to parse user info', e);
  }
}

var updateTodo = function(id, isDone) {
  if (!authToken) {
    showLoginForm();
    return;
  }

  $.ajax({
    url: '/graphql?query=mutation+_{updateTodo(id:"' + id + '",done:' + isDone + ',token:"' + authToken + '"){id,text,done}}'
  }).done(function(data) {
    console.log(data);
    var dataParsed = JSON.parse(data);
    
    if (dataParsed.errors) {
      handleAuthError(dataParsed.errors);
      return;
    }
    
    var updatedTodo = dataParsed.data.updateTodo;
    if (updatedTodo.done) {
      $('#' + updatedTodo.id).parent().parent().parent().addClass('todo-done');
    } else {
      $('#' + updatedTodo.id).parent().parent().parent().removeClass('todo-done');
    }
  });
};

var handleTodoList = function(object) {
  $('.todo-list-container .todo-item').remove();
  var todos = object;

  if (!todos || !todos.length) {
    $('.todo-list-container').append('<p>There are no tasks for you today</p>');
    return;
  } else {
    $('.todo-list-container p').remove();
  }

  $.each(todos, function(i, v) {
    var todoTemplate = $('#todoItemTemplate').html();
    var todo = todoTemplate.replace('{{todo-id}}', v.id);
    todo = todo.replace('{{todo-text}}', v.text);
    todo = todo.replace('{{todo-checked}}', (v.done ? ' checked="checked"' : ''));
    todo = todo.replace('{{todo-done}}', (v.done ? ' todo-done' : ''));

    $('.todo-list-container').append(todo);
    $('#' + v.id).click(function() {
      var id = $(this).prop('id');
      var isDone = $(this).prop('checked');
      updateTodo(id, isDone);
    });
  });
};

var loadTodos = function() {
  if (!authToken) {
    showLoginForm();
    return;
  }

  $.ajax({
    url: "/graphql?query={todoList(token:\"" + authToken + "\"){id,text,done}}"
  }).done(function(data) {
    console.log(data);
    var dataParsed = JSON.parse(data);
    
    if (dataParsed.errors) {
      handleAuthError(dataParsed.errors);
      return;
    }
    
    handleTodoList(dataParsed.data.todoList);
  });
};

var addTodo = function(todoText) {
  if (!authToken) {
    showLoginForm();
    return;
  }

  if (!todoText || todoText === "") {
    alert('Please specify a task');
    return;
  }

  $.ajax({
    url: '/graphql?query=mutation+_{createTodo(text:"' + todoText + '",token:"' + authToken + '"){id,text,done}}'
  }).done(function(data) {
    console.log(data);
    var dataParsed = JSON.parse(data);
    
    if (dataParsed.errors) {
      handleAuthError(dataParsed.errors);
      return;
    }
    
    var todoList = [dataParsed.data.createTodo];
    handleTodoList(todoList);
  });
};

var handleAuthError = function(errors) {
  if (errors && errors.length > 0) {
    // Check for unauthorized error
    var isAuthError = errors.some(function(error) {
      return error.message && error.message.toLowerCase().includes('unauthorized');
    });
    
    if (isAuthError) {
      // Clear stored auth data
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user_info');
      authToken = null;
      currentUser = null;
      
      showLoginForm();
    }
  }
};

var signup = function(email, password) {
  $.ajax({
    url: '/graphql?query=mutation+_{signup(email:"' + email + '",password:"' + password + '"){token,user{id,email}}}'
  }).done(function(data) {
    console.log(data);
    var dataParsed = JSON.parse(data);
    
    if (dataParsed.errors) {
      alert('Signup failed: ' + dataParsed.errors[0].message);
      return;
    }
    
    var result = dataParsed.data.signup;
    authToken = result.token;
    currentUser = result.user;
    
    // Store auth data
    localStorage.setItem('auth_token', authToken);
    localStorage.setItem('user_info', JSON.stringify(currentUser));
    
    // Show todo interface
    showTodoInterface();
    loadTodos();
  });
};

var login = function(email, password) {
  $.ajax({
    url: '/graphql?query=mutation+_{login(email:"' + email + '",password:"' + password + '"){token,user{id,email}}}'
  }).done(function(data) {
    console.log(data);
    var dataParsed = JSON.parse(data);
    
    if (dataParsed.errors) {
      alert('Login failed: ' + dataParsed.errors[0].message);
      return;
    }
    
    var result = dataParsed.data.login;
    authToken = result.token;
    currentUser = result.user;
    
    // Store auth data
    localStorage.setItem('auth_token', authToken);
    localStorage.setItem('user_info', JSON.stringify(currentUser));
    
    // Show todo interface
    showTodoInterface();
    loadTodos();
  });
};

var logout = function() {
  // Clear stored auth data
  localStorage.removeItem('auth_token');
  localStorage.removeItem('user_info');
  authToken = null;
  currentUser = null;
  
  // Show login form
  showLoginForm();
};

var showLoginForm = function() {
  $('.todo-interface').hide();
  $('.auth-forms').show();
  $('.signup-form').hide();
  $('.login-form').show();
};

var showSignupForm = function() {
  $('.todo-interface').hide();
  $('.auth-forms').show();
  $('.login-form').hide();
  $('.signup-form').show();
};

var showTodoInterface = function() {
  $('.auth-forms').hide();
  $('.todo-interface').show();
  
  // Update user info display
  if (currentUser) {
    $('.user-email').text(currentUser.email);
  }
};

$(document).ready(function() {
  // Initialize interface based on auth status
  if (authToken) {
    showTodoInterface();
    loadTodos();
  } else {
    showLoginForm();
  }
  
  // Todo form submission
  $('.todo-add-form').submit(function(e) {
    e.preventDefault();
    addTodo($('.todo-add-form #task').val());
    $('.todo-add-form #task').val('');
  });
  
  // Login form submission
  $('.login-form form').submit(function(e) {
    e.preventDefault();
    var email = $('.login-form #email').val();
    var password = $('.login-form #password').val();
    login(email, password);
  });
  
  // Signup form submission
  $('.signup-form form').submit(function(e) {
    e.preventDefault();
    var email = $('.signup-form #email').val();
    var password = $('.signup-form #password').val();
    signup(email, password);
  });
  
  // Switch to signup form
  $('.show-signup').click(function(e) {
    e.preventDefault();
    showSignupForm();
  });
  
  // Switch to login form
  $('.show-login').click(function(e) {
    e.preventDefault();
    showLoginForm();
  });
  
  // Logout button
  $('.logout-btn').click(function(e) {
    e.preventDefault();
    logout();
  });
});