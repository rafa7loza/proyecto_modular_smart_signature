import { Link } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { ButtonBase } from '@mui/material';
import { MENU_OPEN } from '../../store/actions';

const LogoSection = () => {
  const defaultId = useSelector((state) => state.customization.defaultId);
  const dispatch = useDispatch();
  return (
    <ButtonBase disableRipple onClick={() => dispatch({ type: MENU_OPEN, id: defaultId })} component={Link} to={'/'}>
       <p>Smart Signature</p>
    </ButtonBase>
  );
};

export default LogoSection;
