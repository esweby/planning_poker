import avatar from "animal-avatar-generator";
import styles from "./ListUserItem.module.css";

const ListUserItem = ({ user }: { user: User }) => {
  return (
    <div key={user.username} className={styles.container}>
      <img
        src={`data:image/svg+xml;utf8,${encodeURIComponent(avatar(user?.seed ?? user.username, { size: 40, blackout: false }))}`}
        className={styles.image}
        alt="avatar"
      />
      <h3 className={styles.username}>{user.username}</h3>
      <p className={styles.role}>{user.role}</p>
    </div>
  );
};

export default ListUserItem;
