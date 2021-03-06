import './App.css';
import React, { useState } from 'react';
import Loader from 'react-loader-spinner';

function App() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState([]);
  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage([]);
    setIsLoading(true);

    const formData = {
      name,
      email,
      phone,
      password,
      confirmPassword
    };

    const response = await fetch('/api/mail', {
      method: 'POST',
      body: JSON.stringify(formData),
      headers: {
        'Content-Type': 'application/json'
      }
    });

    let result;
    try {
      result = await response.json();
    } catch (e) {
      console.log('error:', e.message);
    }

    if (result !== undefined) {
      if (result.isOk) {
        setMessage(['Message sent']);
      } else {
        setMessage(result.errorMessage);
      }
    } else {
      setMessage(['Empty response']);
    }
    setIsLoading(false);
  };

  return (
    <div>
      <form className='form' onSubmit={handleSubmit}>
        <div className='main-label'>Registration form</div>
        <div>
          <input
            className='form-input'
            type='text'
            required
            minLength={1}
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder='Name'
          />
        </div>
        <div>
          <input
            className='form-input'
            type='email'
            required
            value={email}
            minLength={5}
            onChange={(e) => setEmail(e.target.value)}
            placeholder='Email'
          />
        </div>
        <div>
          <input
            className='form-input'
            type='text'
            minLength={12}
            maxLength={12}
            required
            value={phone}
            onChange={(e) => {
              setPhone(e.target.value.replace(/[^0-9\d]/gi, ''));
            }}
            placeholder='Phone number'
          />
        </div>
        <div>
          <input
            className='form-input'
            type='password'
            required
            value={password}
            minLength={6}
            maxLength={18}
            onChange={(e) => setPassword(e.target.value)}
            placeholder='Password'
          />
        </div>
        <div>
          <input
            className='form-input'
            type='password'
            required
            value={confirmPassword}
            minLength={6}
            maxLength={18}
            onChange={(e) => setConfirmPassword(e.target.value)}
            placeholder='Confirm password'
          />
        </div>
        <input className='form-input' type='submit' disabled={isLoading} />
        <div className='spinner'>
          {isLoading ? (
            <Loader
              type='Puff'
              color='#00BFFF'
              height={100}
              width={100}
              timeout={30000} //3 secs
            />
          ) : null}
        </div>
        {
          /* eslint-disable react/jsx-key */
          message.map((e) => (
            <div className='error'> {e}</div>
          ))
        }
      </form>
    </div>
  );
}

export default App;
